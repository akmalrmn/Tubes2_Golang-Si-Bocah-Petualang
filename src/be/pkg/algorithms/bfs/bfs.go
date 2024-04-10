package bfs

import (
	"be/pkg/scraper"
	"be/pkg/tree"
	"context"
	"errors"
	"fmt"
	"time"
)

const MaxGoroutines = 12

var HowManyArticleChecked = 0

func processNode(node *tree.Node, nodeCh chan<- *tree.Node, sem chan struct{}) {
	sem <- struct{}{}
	defer func() { <-sem }()
	HowManyArticleChecked++
	links, err := scraper.ExtractLinks("https://en.wikipedia.org" + node.Value)
	if err != nil && !errors.Is(err, context.Canceled) {
		fmt.Println(err)
		return
	}
	for _, link := range links {
		node.AddChild(tree.NewNode(link))
	}
	nodeCh <- node
}

func BidirectionalBreadthFirstSearch(start, end string) []string {
	// ! Profiler !
	startTime := time.Now()
	HowManyArticleChecked = 0
	rootStart, rootEnd := tree.NewNode(start), tree.NewNode(end)
	queueStart, queueEnd := []*tree.Node{rootStart}, []*tree.Node{rootEnd}
	visitedStart, visitedEnd := make(map[string]*tree.Node), make(map[string]*tree.Node)
	nodeChStart, nodeChEnd := make(chan *tree.Node, 10), make(chan *tree.Node, 10)
	sem := make(chan struct{}, MaxGoroutines)
	go processNode(rootStart, nodeChStart, sem)
	go processNode(rootEnd, nodeChEnd, sem)

	for {
		if HowManyArticleChecked%10 == 0 {
			fmt.Println("ArticleChecked ", HowManyArticleChecked, " Time ", time.Since(startTime))
		}
		fmt.Println("ArticlePerSec ", float64(HowManyArticleChecked)/time.Since(startTime).Seconds())

		select {
		case nodeStart := <-nodeChStart:
			visitedStart[nodeStart.Value] = nodeStart
			if nodeEnd, ok := visitedEnd[nodeStart.Value]; ok {
				path := returnPath(nodeStart, nodeEnd)
				if isPathValid(path) {
					printPath(path)
					return path
				}
			}
			processChildren(nodeStart, visitedStart, nodeChStart, sem, queueStart)

		case nodeEnd := <-nodeChEnd:
			visitedEnd[nodeEnd.Value] = nodeEnd
			if nodeStart, ok := visitedStart[nodeEnd.Value]; ok {
				path := returnPath(nodeStart, nodeEnd)
				if isPathValid(path) {
					printPath(path)
					return path
				}
			}
			processChildren(nodeEnd, visitedEnd, nodeChEnd, sem, queueEnd)
		}
	}
}

func processChildren(node *tree.Node, visited map[string]*tree.Node, nodeCh chan<- *tree.Node, sem chan struct{}, queue []*tree.Node) {
	for _, child := range node.Children {
		if _, ok := visited[child.Value]; !ok {
			go processNode(child, nodeCh, sem)
			queue = append(queue, child)
		}
	}
}

func isPathValid(path []string) bool {
	for i := 0; i < len(path)-1; i++ {
		if !linkExists(path[i], path[i+1]) {
			return false
		}
	}
	return true
}

func linkExists(from, to string) bool {
	links, _ := scraper.ExtractLinks(from)
	for _, link := range links {
		if link == to {
			return true
		}
	}
	return false
}

func printPath(path []string) {
	fmt.Println("Path:", path)
}

func returnPath(start, end *tree.Node) []string {
	var path []string
	for n := start; n != nil; n = n.Parent {
		path = append([]string{n.Value}, path...)
	}
	for n := end.Parent; n != nil; n = n.Parent {
		path = append(path, n.Value)
	}
	return path
}

func BidirectionalBreadthFirstSearchOri(start, end string) *tree.Node {
	startTime := time.Now()
	rootStart := tree.NewNode(start)
	rootEnd := tree.NewNode(end)
	queueStart := []*tree.Node{rootStart}
	queueEnd := []*tree.Node{rootEnd}
	visitedStart := make(map[string]*tree.Node)
	visitedEnd := make(map[string]*tree.Node)

	// Create channels to collect nodes that have been processed
	nodeChStart := make(chan *tree.Node)
	nodeChEnd := make(chan *tree.Node)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Create a semaphore to limit the number of goroutines
	sem := make(chan struct{}, MaxGoroutines)

	// Start goroutines for the root nodes
	go processNodeOri(ctx, rootStart, nodeChStart, sem)
	go processNodeOri(ctx, rootEnd, nodeChEnd, sem)

	for {
		if HowManyArticleChecked%10 == 0 {
			fmt.Println("ArticleChecked ", HowManyArticleChecked, " Time ", time.Since(startTime))
		}
		fmt.Println("ArticlePerSec ", float64(HowManyArticleChecked)/time.Since(startTime).Seconds())
		select {
		case nodeStart := <-nodeChStart:
			visitedStart[nodeStart.Value] = nodeStart
			if nodeEnd, ok := visitedEnd[nodeStart.Value]; ok {
				if isPathValid(returnPath(nodeStart, nodeEnd)) {
					return nodeEnd // or some function to reconstruct the path
				}
			}
			for _, child := range nodeStart.Children {
				if _, ok := visitedStart[child.Value]; !ok {
					go processNodeOri(ctx, child, nodeChStart, sem)
					queueStart = append(queueStart, child)
				}
			}

		case nodeEnd := <-nodeChEnd:
			visitedEnd[nodeEnd.Value] = nodeEnd
			if nodeStart, ok := visitedStart[nodeEnd.Value]; ok {
				if isPathValid(returnPath(nodeStart, nodeEnd)) {
					return nodeEnd // or some function to reconstruct the path
				}
			}
			for _, child := range nodeEnd.Children {
				if _, ok := visitedEnd[child.Value]; !ok {
					go processNodeOri(ctx, child, nodeChEnd, sem)
					queueEnd = append(queueEnd, child)
				}
			}
		}
	}
}
func processNodeOri(ctx context.Context, node *tree.Node, nodeCh chan<- *tree.Node, sem chan struct{}) {
	// Acquire a token from the semaphore
	sem <- struct{}{}
	defer func() { <-sem }() // Release the token when the function returns

	// Check the context before making an HTTP request
	if ctx.Err() != nil {
		return
	}

	HowManyArticleChecked++
	links, err := scraper.ExtractLinksAsync(ctx, "https://en.wikipedia.org"+node.Value)

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
