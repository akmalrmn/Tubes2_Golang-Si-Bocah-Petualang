package bfs

import (
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
	_ "sync"
)

type Node struct {
	Value    string
	Children []*Node
	Parent   *Node
}

func NewNode(value string) *Node {
	return &Node{
		Value:    value,
		Children: []*Node{},
		Parent:   nil,
	}
}

func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

func ExtractLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var links []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			if !strings.HasPrefix(href, "#") && strings.HasPrefix(href, "/wiki/") && !strings.Contains(href, ":") {
				links = append(links, href)
			}
		}
	})

	return links, nil
}

func ExtractLinksAsync(ctx context.Context, url string) ([]string, error) {
	result := make(chan []string, 1)
	errResult := make(chan error, 1)

	go func() {
		defer close(result)
		defer close(errResult)
		links, err := ExtractLinks(url)
		if err != nil {
			errResult <- err
			return
		}
		result <- links
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case links := <-result:
		return links, nil
	case err := <-errResult:
		return nil, err
	}
}

const MaxGoroutines = 9999

func BreadthFirstSearch(start, end string) *Node {
	root := NewNode(start)
	queue := []*Node{root}
	visited := make(map[string]bool)

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // make sure all paths cancel the context to avoid context leak

	// Create a channel to collect nodes that have been processed
	nodeCh := make(chan *Node)

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

			fmt.Println("Processing ", node.Value)
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

func processNode(ctx context.Context, node *Node, nodeCh chan<- *Node, sem chan struct{}) {
	// Check if the context has been cancelled before processing the node
	if ctx.Err() != nil {
		return
	}

	// Acquire a token from the semaphore
	sem <- struct{}{}

	links, err := ExtractLinksAsync(ctx, "https://en.wikipedia.org"+node.Value)

	// Only print the error if it's not a context cancellation
	if err != nil && !errors.Is(err, context.Canceled) {
		fmt.Println(err)
		return
	}

	for _, link := range links {
		child := NewNode(link)
		node.AddChild(child)
	}
	defer func() { <-sem }() // Release the token when the function returns
	// Send the processed node to the channel
	nodeCh <- node
}
