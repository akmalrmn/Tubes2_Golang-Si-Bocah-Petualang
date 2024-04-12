package bfs

import (
	"be/pkg/scraper"
	"be/pkg/tree"
	"context"
	"errors"
	"fmt"
	"github.com/korovkin/limiter"
	"sync"
	"time"
	"unsafe"
)

type JobInfo struct {
	Links string
	Err   error
}

const MaxGoroutines = 100

var (
	HowManyArticleChecked = 0
	mutex                 = &sync.Mutex{}
)

func BidirectionalBreadthFirstSearch(start, end string) []string {
	// ! Profiler !
	startTime := time.Now()
	limit := limiter.NewConcurrencyLimiter(MaxGoroutines)

	HowManyArticleChecked = 0
	rootStart, rootEnd := tree.NewNode(start), tree.NewNode(end)
	queueStart, queueEnd := []*tree.Node{rootStart}, []*tree.Node{rootEnd}
	visitedStart, visitedEnd := make(map[string]*tree.Node), make(map[string]*tree.Node)
	nodeChStart, nodeChEnd := make(chan *tree.Node, 9999), make(chan *tree.Node, 9999)

	go processNode(rootStart, nodeChStart, limit)
	go processNode(rootEnd, nodeChEnd, limit)

	for {
		fmt.Println("Size Of Ch : ", unsafe.Sizeof(nodeChStart), unsafe.Sizeof(nodeChEnd))
		mutex.Lock()
		if HowManyArticleChecked%10 == 0 {
			fmt.Println("ArticleChecked ", HowManyArticleChecked, " Time ", time.Since(startTime))
		}
		fmt.Println("ArticlePerSec ", float64(HowManyArticleChecked)/time.Since(startTime).Seconds())
		mutex.Unlock()
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
			processChildren(nodeStart, visitedStart, nodeChStart, limit, queueStart)

		case nodeEnd := <-nodeChEnd:
			visitedEnd[nodeEnd.Value] = nodeEnd
			if nodeStart, ok := visitedStart[nodeEnd.Value]; ok {
				path := returnPath(nodeStart, nodeEnd)
				if isPathValid(path) {
					printPath(path)
					return path
				}
			}
			processChildren(nodeEnd, visitedEnd, nodeChEnd, limit, queueEnd)
		}
	}
}

func processNode(node *tree.Node, nodeCh chan<- *tree.Node, limit *limiter.ConcurrencyLimiter) {
	_, err := limit.Execute(func() {
		mutex.Lock()
		HowManyArticleChecked++
		mutex.Unlock()
		links, err := scraper.ExtractLinks("https://en.wikipedia.org" + node.Value)
		if err != nil && !errors.Is(err, context.Canceled) {
			return
		}
		for _, link := range links {
			node.AddChild(tree.NewNode(link))
		}
		nodeCh <- node
	})
	if err != nil {
		return
	}
}

func processChildren(node *tree.Node, visited map[string]*tree.Node, nodeCh chan<- *tree.Node, limit *limiter.ConcurrencyLimiter, queue []*tree.Node) {
	for _, child := range node.Children {
		if _, ok := visited[child.Value]; !ok {
			processNode(child, nodeCh, limit)
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
