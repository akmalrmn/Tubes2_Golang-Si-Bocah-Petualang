package ids

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"log"
	"os"
	"github.com/gocolly/colly"
)

type Result struct {
	Path          []string
	Degrees       int
	TimeTaken     time.Duration
	LinksVisited  int
}

func IterativeDeepeningSearch(start, ends string) *Result{
    start = "https://en.wikipedia.org/" + start
    ends = "https://en.wikipedia.org/" + ends
	c := colly.NewCollector(
    colly.Async(true),
    )

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*wikipedia.*",
		Parallelism: 15,
		Delay:       100 * time.Millisecond,
	})

	// The maximum depth for the IDS
	maxDepth := 6

	var result Result

	// Create a map to store the depth of each URL
	depths := make(map[string]int)

	// Mutex to protect the depths map
	var depthsMutex sync.RWMutex

	targetUrl := ends

	// Variable to store the least depth at which the target URL was found
	leastDepth := maxDepth + 1

	// WaitGroup to wait for all URLs to be visited
	var wg sync.WaitGroup

	// Counter for the number of links visited
	var linksVisited int
	parents := make(map[string]string)
	// Start the timer
	startTime := time.Now()
	ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
		go func() {
			for range ticker.C {
					fmt.Println("Number of links visited:", linksVisited)
					fmt.Println("Time taken:", time.Since(startTime))
			}
	}()

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if isValidLink(link) {
			url := e.Request.AbsoluteURL(link)
			depthsMutex.RLock()
			currentDepth, ok := depths[e.Request.URL.String()]
			depthsMutex.RUnlock()
			if !ok {
				return
			}
			if currentDepth < maxDepth {
				depthsMutex.Lock()
				oldDepth, found := depths[url]
				if !found || oldDepth > currentDepth+1 {
					depths[url] = currentDepth + 1
					parents[url] = e.Request.URL.String()
					depthsMutex.Unlock()
					if url == targetUrl && currentDepth+1 < leastDepth {
						leastDepth = currentDepth + 1
						result.Degrees = leastDepth
						fmt.Println(url)
						fmt.Println("\nFound target URL at depth", leastDepth)
						fmt.Println("Time taken:", time.Since(startTime))
            fmt.Println("Number of links visited:", linksVisited)
						path := []string{url}
                for url != start {
                    url = parents[url]
                    path = append([]string{url}, path...)
                }
                fmt.Println("Path:", strings.Join(path, " -> "))
								result.Path = path
								result.TimeTaken = time.Since(startTime)
								result.LinksVisited = linksVisited
						os.Exit(0)
					}
					wg.Add(1)
        go func() {
					defer wg.Done()
            err := c.Visit(url)
            if err != nil {
							log.Println("Error visiting URL:", url, "Error:", err)
					} 
        }()
        linksVisited++
				} else {
					depthsMutex.Unlock()
				}
			}
		}
	})

	// Initialize the depths map with the starting URL
	depthsMutex.Lock()
	depths[start] = 0
	depthsMutex.Unlock()
	wg.Add(1)
    go func() {
			defer wg.Done()
        c.Visit(start)
        
    }()

	// Wait for all URLs to be visited
wg.Wait()
c.Wait()
				log.Println("adadadadaadad")
	// Print the least depth at which the target URL was found
	if leastDepth > maxDepth {
    fmt.Println("The target URL was not found within the maximum depth")
    fmt.Println("Number of links visited:", linksVisited)
    fmt.Println("Time taken:", time.Since(startTime))
    os.Exit(1)  // Exit with a status indicating the target was not found
} else {
    fmt.Println("The least depth at which the target URL was found is", leastDepth)
    fmt.Println("Number of links visited:", linksVisited)
    fmt.Println("Time taken:", time.Since(startTime))
    os.Exit(0)  // Exit successfully as the target URL was found
}

	// Print the number of links visited and the time taken
	fmt.Println("Number of links visited:", linksVisited)
	fmt.Println("Time taken:", time.Since(startTime))

    return &result

}

func isValidLink(link string) bool {
	prefixes := []string{
		"/wiki/Main_Page",
		"/wiki/File",
		"/wiki/Special",
		"/wiki/Wikipedia",
		"/wiki/Help",
		"/wiki/Portal",
		"/wiki/Template",
		"/wiki/Category",
		"/wiki/Talk",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(link, prefix) {
			return false
		}
	}
	return strings.HasPrefix(link, "/wiki")
}