package bfs

import (
	"be/pkg/scraper"
	"be/pkg/set"
	"be/pkg/tree"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"runtime"
	"sync"
	"time"
)

///  * ============== Var =============== * ///

var (
	maxGoRoutines = 10000

	chanOutForward  = make(chan *tree.Node, 1000)
	chanOutBackward = make(chan *tree.Node, 1000)

	tasksForward  = make(chan *tree.Node, 1000)
	tasksBackward = make(chan *tree.Node, 1000)

	visitedForward  sync.Map
	visitedBackward sync.Map

	pathResult = set.NewSetOfSlice()

	sem = make(chan struct{}, maxGoRoutines)

	finished      = false
	increaseLimit = 30000
)

func reset() {
	chanOutForward = make(chan *tree.Node, 1000)
	chanOutBackward = make(chan *tree.Node, 1000)
	visitedForward = sync.Map{}
	visitedBackward = sync.Map{}
	pathResult = set.NewSetOfSlice()
	finished = false
	maxGoRoutines = 10000
}

///  * ============== Function =============== * ///

func SearchWithTimeout(start, end string, timeout time.Duration) [][]string {
	// Create a context that will automatically cancel after the specified timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Call BidirectionalBreadthFirstSearch with the context
	result := BidirectionalBreadthFirstSearch(ctx, start, end)

	// If the result is empty and the context has not been cancelled, continue the search
	for len(result) == 0 && ctx.Err() == nil {
		result = BidirectionalBreadthFirstSearch(ctx, start, end)
	}

	return result
}

func BidirectionalBreadthFirstSearch(ctx context.Context, start, end string) [][]string {
	go trackGoroutines() // ! Profiler ! //
	defer func() { finished = true }()
	reset()

	rootStart := tree.NewNode(start)
	rootEnd := tree.NewNode(end)

	// Create channels to collect nodes that have been processed
	nodeChStart := make(chan *tree.Node, 99999)
	nodeChEnd := make(chan *tree.Node, 99999)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Create a semaphore to limit the number of goroutines

	go func() { // Increase the limit of goroutines every 5 seconds
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if !finished && maxGoRoutines < increaseLimit {
					maxGoRoutines += 100 // Increase the limit by 10
				}
			}
		}
	}()

	go processNode(ctx, rootStart, nodeChStart)
	go processNode(ctx, rootEnd, nodeChEnd)

	// Visited start and end nodes
	visitedForward.Store(start, rootStart)
	visitedBackward.Store(end, rootEnd)

	for {
		select {
		case <-ctx.Done():
			return pathResult.ToSlice()
		case nodeStart := <-nodeChStart:
			visitedForward.Store(nodeStart.Value, nodeStart)
			if nodes, ok := visitedBackward.Load(nodeStart.Value); ok {
				path := returnPathBiBFS(nodeStart, nodes.(*tree.Node))
				go checkPathValidity(path)
			}

			for _, child := range nodeStart.Children {
				if _, ok := visitedForward.Load(child.Value); !ok {
					visitedBackward.Store(child.Value, child)
					if child.Value == end {
						path := returnPathBFS(child)
						go checkPathValidity(path)
					}
					sem <- struct{}{} // This will block if sem is full
					go func() {
						processNode(ctx, child, nodeChStart)
					}()
				}
			}

		case nodeEnd := <-nodeChEnd:
			visitedBackward.Store(nodeEnd.Value, nodeEnd)
			if nodes, ok := visitedForward.Load(nodeEnd.Value); ok {
				path := returnPathBiBFS(nodes.(*tree.Node), nodeEnd)
				go checkPathValidity(path)
			}
			for _, child := range nodeEnd.Children {
				if _, ok := visitedBackward.Load(child.Value); !ok {
					visitedBackward.Store(child.Value, child)
					if child.Value == start {
						path := returnPathBFS(child)
						path = reverse(path)
						go checkPathValidity(path)
					}
					sem <- struct{}{} // This will block if sem is full
					go func() {
						processNode(ctx, child, nodeChEnd)
					}()
				}
			}
		}
	}
}

func processNode(ctx context.Context, node *tree.Node, nodeCh chan<- *tree.Node) {
	// Check the context before making an HTTP request
	if ctx.Err() != nil {
		return
	}
	links, err := scraper.ExtractLinksAsync("https://en.wikipedia.org" + node.Value)

	defer func() { <-sem }()

	// Only print the error if it's not a context cancellation
	if err != nil && !errors.Is(err, context.Canceled) {
		fmt.Println(err)
		return
	}

	for _, link := range links {
		child := tree.NewNode(link)
		node.AddChild(child)
	}

	// Send the processed node to the channel
	nodeCh <- node
}

// BiDirectionalBFS
//   - BiDirectionalBFS is a function that finds the shortest path between two nodes using the Bi-Directional Breadth First Search algorithm.
//   - It takes two parameters, start and end, which are the start and end nodes respectively.
//   - It returns a slice of strings that represents the shortest path between the two nodes.
func BiDirectionalBFS(start, end string) [][]string {
	go trackGoroutines() // ! Profiler ! //
	defer func() { finished = true }()
	reset()

	timeStart := time.Now()

	// Delete Links Cache
	scraper.LinkCache.Purge()

	// Create start and end nodes
	startNode := tree.NewNode(start)
	endNode := tree.NewNode(end)

	// Create worker pool
	for i := 0; i < maxGoRoutines/2; i++ {
		go workerBi(tasksForward, 1)
		go workerBi(tasksBackward, 2)
	}

	// Process the start and end nodes
	tasksForward <- startNode
	tasksBackward <- endNode

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
		timeNow := time.Since(timeStart)
		if timeNow.Seconds() > 10 && pathResult.Size() > 0 {
			return pathResult.ToSlice()
		}
		select {
		case node := <-chanOutForward:
			if nodes, ok := visitedBackward.Load(node.Value); ok {
				path := returnPathBiBFS(nodes.(*tree.Node), node)
				go checkPathValidity(path)
			}
		case node := <-chanOutBackward:
			if nodes, ok := visitedForward.Load(node.Value); ok {
				path := returnPathBiBFS(nodes.(*tree.Node), node)
				go checkPathValidity(path)
			}
		}
	}
}

///  * ============== Helper =============== * ///

func checkPathValidity(path []string) {
	results := make([]<-chan bool, len(path)-1)

	for i := 0; i < len(path)-1; i++ {
		results[i] = checkLinkContainAsync(path[i], path[i+1])
	}

	for _, result := range results {
		if !<-result {
			fmt.Println("Path ", path, " is invalid")
			return
		}
	}
	fmt.Println("Path ", path, " is valid")
	pathResult.Add(path)
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

func returnPathBFS(node *tree.Node) []string {
	var path []string
	for node != nil {
		path = append(path, node.Value)
		node = node.Parent
	}
	return reverse(path)
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

func workerBi(tasks <-chan *tree.Node, code int) {
	for {
		select {
		case node, _ := <-tasks:
			if finished {
				return
			}
			if code == 1 {
				ProcessNodeBi(node, chanOutForward)
			} else {
				ProcessNodeBi(node, chanOutBackward)
			}
		}
	}
}

func trackGoroutines() {
	start := time.Now()
	for {
		time.Sleep(time.Second * 5)
		fmt.Println(" ======================== ")
		fmt.Println("Article per sec : ", scraper.NumOfArticlesProcessed/int(time.Since(start).Seconds()))
		fmt.Println("Number of articles processed: ", scraper.NumOfArticlesProcessed)
		fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
		fmt.Println("Max number of goroutines: ", maxGoRoutines)
		fmt.Println("Semaphore :", len(sem))
	}
}
