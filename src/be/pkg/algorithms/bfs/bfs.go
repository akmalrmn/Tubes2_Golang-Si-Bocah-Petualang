package bfs

import (
	"be/pkg/tree"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func BFS(starts, ends string) []string {

	client := &http.Client{
		Timeout: 30 * time.Second, // Set timeout to 30 seconds
	}

	// Instantiate a new collector
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0"),
		colly.URLFilters(
			regexp.MustCompile(`https://en\.wikipedia\.org/wiki/.*`),
		),
		colly.WithTransport(client.Transport), // Use the custom HTTP client
	)

	// Create a request queue with 16 consumer threads
	q, _ := queue.New(
		16, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	// Limit the number of threads started by colly
	c.Limit(&colly.LimitRule{
		DomainGlob:  "en.wikipedia.org",
		Parallelism: 400,
	})

	// Print the URL of the page visited
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Print error if something went wrong
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	// Make a channel to receive the result
	result := make(chan []string)

	// Set the URL to start scraping from
	q.AddURL("https://en.wikipedia.org" + starts)

	// Root Tree
	root := tree.NewNode(starts)

	// Set the handler for the start URL
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if !strings.HasPrefix(link, "/wiki/Main_Page") && !strings.HasPrefix(link, "/wiki/File") && !strings.HasPrefix(link, "/wiki/Special") && !strings.HasPrefix(link, "/wiki/Wikipedia") && !strings.HasPrefix(link, "/wiki/Help") && !strings.HasPrefix(link, "/wiki/Portal") && !strings.HasPrefix(link, "/wiki/Template") && !strings.HasPrefix(link, "/wiki/Category") && !strings.HasPrefix(link, "/wiki/Talk") && !strings.HasPrefix(link, "/wiki/Wikipedia") {
			q.AddURL(e.Request.AbsoluteURL(link))
			node := tree.NewNode(link)
			root.AddChild(node)

			if link == ends {
				result <- traverseToRoot(node)
			}
		}
	})

	for {
		q.Run(c)
	}

	// Return the result
	return <-result
}

func traverseToRoot(endNode *tree.Node) []string {
	var path []string
	for endNode != nil {
		path = append(path, endNode.Value)
		endNode = endNode.Parent
	}
	return path
}
