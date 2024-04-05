package bfs

import (
	"context"
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

func BreadthFirstSearch(start, end string) *Node {
	root := NewNode(start)
	queue := []*Node{root}
	visited := make(map[string]bool)

	// Create a channel to collect nodes that have been processed
	nodeCh := make(chan *Node)

	// Start a goroutine for the root node
	go processNode(root, nodeCh)

	for node := range nodeCh {
		visited[node.Value] = true

		if node.Value == end {
			// Trace back the path from the destination to the source
			path := []string{}
			for n := node; n != nil; n = n.Parent {
				path = append([]string{n.Value}, path...)
			}
			fmt.Println("Path:", path)
			return node
		}

		for _, child := range node.Children {
			if !visited[child.Value] {
				go processNode(child, nodeCh)
				queue = append(queue, child)
			}
		}
	}

	return nil
}

func processNode(node *Node, nodeCh chan<- *Node) {
	fmt.Println("Processing node:", node.Value) // Print the node value
	links, err := ExtractLinksAsync(context.Background(), "https://en.wikipedia.org"+node.Value)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, link := range links {
		child := NewNode(link)
		node.AddChild(child)
	}
	// Send the processed node to the channel
	nodeCh <- node
}
