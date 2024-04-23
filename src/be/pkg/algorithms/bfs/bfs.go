package bfs

import (
	priorityqueue "be/pkg/PriorityQueue"
	"be/pkg/scraper"
	"be/pkg/set"
	"be/pkg/tree"
	"fmt"
	"runtime"
	"sync"
	"time"
)

///  * ============== Var =============== * ///

var (
	maxGoRoutines = 10 // Maksismal jumlah goroutine yang akan dijalankan

	// Channel untuk mengirimkan node yang sudah di proses
	chanOutForward  = make(chan *tree.Node, 10000)
	chanOutBackward = make(chan *tree.Node, 9999)

	// Map untuk menyimpan node yang sudah di proses
	visitedForward  sync.Map
	visitedBackward sync.Map

	// Set untuk menyimpan hasil path && status apakah sudah selesai
	finished   = false
	pathResult = set.NewSetOfSlice()

	// Max limit untuk goroutine yang akan ditingkatkan
	increaseLimit = 800

	// TreeRoot Node root untuk tree awal debugging
	TreeRoot *tree.Node
)

/*
 * Reset function, digunakan untuk mereset semua variabel global
 */
func reset() {
	maxGoRoutines = 10
	chanOutForward = make(chan *tree.Node, 10000)
	chanOutBackward = make(chan *tree.Node, 10000)
	visitedForward = sync.Map{}
	visitedBackward = sync.Map{}
	pathResult = set.NewSetOfSlice()
	finished = false
	increaseLimit = 800
	TreeRoot = nil
}

///  * ============== Function =============== * ///

// BiDirectionalBFS
//   - BiDirectionalBFS is a function that finds the shortest path between two nodes using the Bi-Directional Breadth First Search algorithm.
//   - It takes two parameters, start and end, which are the start and end nodes respectively.
//   - It returns a slice of strings that represents the shortest path between the two nodes.
func BiDirectionalBFS(start, end string) {
	defer func() { finished = true }() // Set finished to true when the function returns
	reset()                            // Reset all global variables

	startTime := time.Now()

	// Delete Links Cache
	scraper.LinkCache.Purge()

	// Create start and end nodes
	startNode := tree.NewNode(start)
	TreeRoot = startNode
	endNode := tree.NewNode(end)

	// Create priority queue
	tasksForward := priorityqueue.NewPriorityChannel()
	tasksBackward := priorityqueue.NewPriorityChannel()

	// Create worker pool
	for i := 0; i < maxGoRoutines/2; i++ {
		go workerBi(tasksForward, 1)
		go workerBi(tasksBackward, 2)
	}

	// Add visitedForward and visitedBackward
	visitedForward.Store(start, startNode)
	visitedBackward.Store(end, endNode)

	// Process the start and end nodes
	tasksForward.Add(startNode)
	tasksBackward.Add(endNode)

	go func() { // Increase the limit of goroutines every 5 seconds
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if !finished && maxGoRoutines < increaseLimit {
					maxGoRoutines += 10 // Increase the limit by 10
					for i := 0; i < 5; i++ {
						go workerBi(tasksForward, 1)
						go workerBi(tasksBackward, 2)
					}
				}
			}
		}
	}()

	for {
		if time.Since(startTime).Minutes() > 5 && pathResult.Size() > 0 { // If the execution time exceeds 1 minute and there is a path result
			return
		}
		select {
		/*
			Process the output of the workers
			if the node is found in the other direction, check the path
			if the node is not visited, add it to the tasks
		*/
		case val := <-chanOutForward:
			// Only access visitedBackward here
			if node, ok := visitedBackward.Load(val.Value); ok {
				path := returnPathBiBFS(val, node.(*tree.Node))
				go checkPathValidity(path)
			}
			if _, ok := visitedForward.Load(val.Value); !ok {
				tasksForward.Add(val)
			}

		/*
			Process the tasks
			set the node as visited
			process the node
		*/
		case tasks := <-tasksForward.C():
			// Only access visitedForward here
			visitedForward.Store(tasks.Value.Value, tasks.Value)
			go ProcessNodeBi(tasks.Value, chanOutForward)

		/*
			Process the output of the workers
			if the node is found in the other direction, check the path
			if the node is not visited, add it to the tasks
		*/
		case val := <-chanOutBackward:
			// Only access visitedForward here
			if node, ok := visitedForward.Load(val.Value); ok {
				path := returnPathBiBFS(node.(*tree.Node), val)
				go checkPathValidity(path)
			}
			if _, ok := visitedBackward.Load(val.Value); !ok {
				tasksBackward.Add(val)
			}

		/*
			Process the tasks
			set the node as visited
			process the node
		*/
		case tasks := <-tasksBackward.C():
			visitedBackward.Store(tasks.Value.Value, tasks.Value)
			go ProcessNodeBi(tasks.Value, chanOutBackward)
		}
	}
}

///  * ============== Helper =============== * ///

// GetPathResult is a function that returns the result of the BiDirectionalBFS function.
func GetPathResult(start, end string) [][]string {
	startTime := time.Now()
	go BiDirectionalBFS(start, end)
	go trackGoroutines() // ! Profiler ! //
	for {
		if pathResult.Size() > 0 && time.Since(startTime).Minutes() > 5 {
			return pathResult.ToSlice()
		}
	}
}

// checkPathValidity is a function that checks the validity of the path.
func checkPathValidity(path []string) {
	results := make([]<-chan bool, len(path)-1)

	for i := 0; i < len(path)-1; i++ {
		results[i] = checkLinkContainAsync(path[i], path[i+1])
	}

	for _, result := range results {
		if !<-result {
			return
		}
	}

	pathResult.Add(path)
}

// checkLinkContainAsync is a function that checks if a link contains another link asynchronously.
func checkLinkContainAsync(start, end string) <-chan bool {
	result := make(chan bool)
	go func() {
		defer close(result)
		result <- linkContain(start, end)
	}()
	return result
}

// linkContain is a function that checks if a link contains another link.
func linkContain(start, end string) bool {
	links, _ := scraper.ExtractLinksNonAdd("https://en.wikipedia.org" + start)

	for _, link := range links {
		if link == end {
			return true
		}
	}
	return false
}

// returnPathBiBFS is a function that returns the path of the Bi-Directional BFS algorithm.
func returnPathBiBFS(midStart, midEnd *tree.Node) []string {
	var path []string
	for midStart != nil {
		path = append(path, midStart.Value)
		midStart = midStart.Parent
	}
	path = reverse(path)
	midEnd = midEnd.Parent
	for midEnd != nil {
		path = append(path, midEnd.Value)
		midEnd = midEnd.Parent
	}
	return path
}

// reverse is a function that reverses a slice of strings.
func reverse(path []string) []string {
	for i := 0; i < len(path)/2; i++ {
		j := len(path) - i - 1
		path[i], path[j] = path[j], path[i]
	}
	return path
}

///  * ============== Process Node =============== * ///

// ProcessNodeBi is a function that processes a node in the Bi-Directional BFS algorithm.
func ProcessNodeBi(node *tree.Node, output chan<- *tree.Node) {
	links, err := scraper.ExtractLinksAsync("https://en.wikipedia.org" + node.Value)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, link := range links {
		if link == node.Value {
			continue
		}
		child := *tree.NewNode(link)
		node.AddChild(&child)
		output <- &child
	}
}

/// * ============== Worker =============== * ///

// workerBi is a function that processes the nodes in the Bi-Directional BFS algorithm.
func workerBi(tasks *priorityqueue.PriorityChannel, code int) {
	for node := range tasks.C() {
		if code == 1 {
			ProcessNodeBi(node.Value, chanOutForward)
		} else {
			ProcessNodeBi(node.Value, chanOutBackward)
		}
	}
}

// / ! ============== Profiler =============== ! ///
// trackGoroutines is a function that tracks the number of goroutines, the number of articles processed, and the maximum number of goroutines.
func trackGoroutines() {
	start := time.Now()
	for {
		time.Sleep(time.Second * 5) // update every 2 seconds
		fmt.Println(" ======================== ")
		fmt.Println("Article per sec : ", scraper.NumOfArticlesProcessed/int(time.Since(start).Seconds()))
		fmt.Println("Number of articles processed: ", scraper.NumOfArticlesProcessed)
		fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
		fmt.Println("Max number of goroutines: ", maxGoRoutines)
		fmt.Println(" Path Result: ", pathResult.Size())
	}
}
