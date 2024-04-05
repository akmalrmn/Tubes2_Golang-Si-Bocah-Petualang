package bfs

import (
	"be/pkg/scraper"
	"be/pkg/tree"
	"context"
	"errors"
	"fmt"
)

const MaxGoroutines = 50 // Lower the limit of goroutines

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

	for {
		select {
		case node := <-nodeCh:
			visited[node.Value] = true

			if node.Value == end {
				// Trace back the path from the destination to the source
				var path []string
				for n := node; n != nil; n = n.Parent {
					path = append([]string{n.Value}, path...)
				}
				fmt.Println("Path:", path)
				return node
			}

			for _, child := range node.Children {
				if !visited[child.Value] {
					go processNode(ctx, child, nodeCh, sem)
					queue = append(queue, child)
				}
			}
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
