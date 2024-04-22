package bfs

import (
	priorityqueue "be/pkg/PriorityQueue"
	"be/pkg/scraper"
	"be/pkg/tree"
	"fmt"
	"runtime"
	"sync"
	"time"
)

///  * ============== Var =============== * ///

var (
	maxGoRoutines = 10
	chanOut       = make(chan *tree.Node, 1000)

	chanOutForward  = make(chan *tree.Node, 1000)
	chanOutBackward = make(chan *tree.Node, 1000)

	visitedForward  sync.Map
	visitedBackward sync.Map

	finished      = false
	increaseLimit = 600
)

///  * ============== Function =============== * ///

// BiDirectionalBFS
//   - BiDirectionalBFS is a function that finds the shortest path between two nodes using the Bi-Directional Breadth First Search algorithm.
//   - It takes two parameters, start and end, which are the start and end nodes respectively.
//   - It returns a slice of strings that represents the shortest path between the two nodes.
func BiDirectionalBFS(start, end string) []string {
	go trackGoroutines() // ! Profiler ! //
	defer func() { finished = true }()

	// Delete Links Cache
	scraper.LinkCache.Purge()

	// Create start and end nodes
	startNode := tree.NewNode(start)
	endNode := tree.NewNode(end)

	tasksForward := priorityqueue.NewPriorityChannel()
	tasksBackward := priorityqueue.NewPriorityChannel()

	// Create worker pool
	for i := 0; i < maxGoRoutines/2; i++ {
		go workerBi(tasksForward, 1)
		go workerBi(tasksBackward, 2)
	}

	// Process the start and end nodes
	tasksForward.Add(startNode)
	tasksBackward.Add(endNode)

	// Visited start and end nodes
	visitedForward.Store(start, startNode)
	visitedBackward.Store(end, endNode)

	go func() {
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
		select {
		case val := <-chanOutForward:
			// Only access visitedBackward here
			if node, ok := visitedBackward.Load(val.Value); ok {
				path := returnPathBiBFS(val, node.(*tree.Node))
				if checkPathValidity(path) {
					return path
				}
			}
			// Only access visitedForward here
			if _, ok := visitedForward.Load(val.Value); !ok {
				visitedForward.Store(val.Value, val)
				tasksForward.Add(val)
			}
		case tasks := <-tasksForward.C():
			// Only access visitedForward here
			if _, ok := visitedForward.Load(tasks.Value.Value); !ok {
				visitedForward.Store(tasks.Value.Value, tasks)
				go ProcessNodeBi(tasks.Value, chanOutForward)
			}
		case val := <-chanOutBackward:
			// Only access visitedForward here
			if node, ok := visitedForward.Load(val.Value); ok {
				path := returnPathBiBFS(node.(*tree.Node), val)
				if checkPathValidity(path) {
					return path
				}
			}
			// Only access visitedBackward here
			if _, ok := visitedBackward.Load(val.Value); !ok {
				visitedBackward.Store(val.Value, val)
				tasksBackward.Add(val)
			}
		case tasks := <-tasksBackward.C():
			// Only access visitedBackward here
			if _, ok := visitedBackward.Load(tasks.Value.Value); !ok {
				visitedBackward.Store(tasks.Value.Value, tasks)
				go ProcessNodeBi(tasks.Value, chanOutBackward)
			}
		}
	}
}

///  * ============== Helper =============== * ///

func checkPathValidity(path []string) bool {
	results := make([]<-chan bool, len(path)-1)

	for i := 0; i < len(path)-1; i++ {
		results[i] = checkLinkContainAsync(path[i], path[i+1])
	}

	for _, result := range results {
		if !<-result {
			return false
		}
	}

	return true
}

func checkLinkContainAsync(start, end string) <-chan bool {
	result := make(chan bool)
	go func() {
		defer close(result)
		result <- linkContain(start, end)
	}()
	return result
}

func linkContain(start, end string) bool {
	links, _ := scraper.ExtractLinksNonAdd("https://en.wikipedia.org" + start)

	for _, link := range links {
		if link == end {
			return true
		}
	}
	return false
}

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

func reverse(path []string) []string {
	for i := 0; i < len(path)/2; i++ {
		j := len(path) - i - 1
		path[i], path[j] = path[j], path[i]
	}
	return path
}

///  * ============== Process Node =============== * ///

func ProcessNodeBi(node *tree.Node, output chan *tree.Node) {
	links, err := scraper.ExtractLinksAsync("https://en.wikipedia.org" + node.Value)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, link := range links {
		child := *tree.NewNode(link)
		node.AddChild(&child)
		output <- &child
	}
}

/// * ============== Worker =============== * ///

func workerBi(tasks *priorityqueue.PriorityChannel, code int) {
	for node := range tasks.C() {
		if code == 1 {
			ProcessNodeBi(node.Value, chanOutForward)
		} else {
			ProcessNodeBi(node.Value, chanOutBackward)
		}
	}
}

func trackGoroutines() {
	start := time.Now()
	for {
		time.Sleep(time.Second * 5) // update every 2 seconds
		fmt.Println(" ======================== ")
		fmt.Println("Article per sec : ", scraper.NumOfArticlesProcessed/int(time.Since(start).Seconds()))
		fmt.Println("Number of articles processed: ", scraper.NumOfArticlesProcessed)
		fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
		fmt.Println("Max number of goroutines: ", maxGoRoutines)
	}
}
