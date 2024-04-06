package bfs

import (
	"be/pkg/scraper"
	"be/pkg/tree"
	"context"
	"errors"
	"fmt"
)

const MaxGoroutines = 12

var HowManyArticelChecked = 0

func BreadthFirstSearch(start, end string) *tree.Node {
	root := tree.NewNode(start)
	queue := []*tree.Node{root}
	visited := make(map[string]bool)

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // make sure all paths cancel the context to avoid context leak

	// Create a channel to collect nodes that have been processed
	nodeCh := make(chan *tree.Node)

	// Create a semaphore to limit the number of goroutines
	sem := make(chan struct{}, MaxGoroutines)

	// Start a goroutine for the root node
	go processNode(ctx, root, nodeCh, sem)

	for { // Loop until the queue is empty
		select { // Select the first available case
		case node := <-nodeCh: // If a node is processed
			visited[node.Value] = true
			fmt.Println("Visited:", node.Value, " , ", node.Value == end)
			if node.Value == end { // If the destination is found
				// Trace back the path from the destination to the source
				var path []string
				for n := node; n != nil; n = n.Parent {
					path = append([]string{n.Value}, path...)
				}
				fmt.Println("Path:", path)
				return node
			}

			for _, child := range node.Children { // Add unvisited children to the queue
				if !visited[child.Value] { // If the child is not visited
					go processNode(ctx, child, nodeCh, sem) // Start a goroutine for the child
					queue = append(queue, child)
				}
			}
			// Remove processed nodes from the queue
			queue = queue[1:]
		default:
			// If there's no node to process, return nil
			if len(queue) == 0 {
				return nil
			}
		}
	}
}

func processNode(ctx context.Context, node *tree.Node, nodeCh chan<- *tree.Node, sem chan struct{}) {
	// Acquire a token from the semaphore
	sem <- struct{}{}
	defer func() { <-sem }() // Release the token when the function returns

	// Check the context before making an HTTP request
	if ctx.Err() != nil {
		return
	}

	HowManyArticelChecked++
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

func BidirectionalBreadthFirstSearch(start, end string) *tree.Node {
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
	go processNode(ctx, rootStart, nodeChStart, sem)
	go processNode(ctx, rootEnd, nodeChEnd, sem)

	for {
		select {
		case nodeStart := <-nodeChStart:
			visitedStart[nodeStart.Value] = nodeStart
			if nodeEnd, ok := visitedEnd[nodeStart.Value]; ok {
				printPath(nodeStart, nodeEnd)
				return nodeStart // or some function to reconstruct the path
			}
			for _, child := range nodeStart.Children {
				if _, ok := visitedStart[child.Value]; !ok {
					go processNode(ctx, child, nodeChStart, sem)
					queueStart = append(queueStart, child)
				}
			}

		case nodeEnd := <-nodeChEnd:
			visitedEnd[nodeEnd.Value] = nodeEnd
			if nodeStart, ok := visitedStart[nodeEnd.Value]; ok {
				printPath(nodeStart, nodeEnd)
				return nodeEnd // or some function to reconstruct the path
			}
			for _, child := range nodeEnd.Children {
				if _, ok := visitedEnd[child.Value]; !ok {
					go processNode(ctx, child, nodeChEnd, sem)
					queueEnd = append(queueEnd, child)
				}
			}
		}
	}
}

func printPath(start, end *tree.Node) {
	var path []string
	for n := start; n != nil; n = n.Parent {
		path = append([]string{n.Value}, path...)
	}
	for n := end.Parent; n != nil; n = n.Parent {
		path = append(path, n.Value)
	}
	fmt.Println("Path:", path)
}
