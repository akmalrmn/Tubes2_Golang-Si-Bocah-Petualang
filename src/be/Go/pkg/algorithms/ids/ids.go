package ids

import (
	"be/pkg/Result"
	"be/pkg/set"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"github.com/gocolly/colly/v2"
)

type contextTransport struct {
	ctx   context.Context
	trans *http.Transport
}

func (t *contextTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.WithContext(t.ctx)
	return t.trans.RoundTrip(req)
}

func collectorWithContext(c *colly.Collector, ctx context.Context) {
	c.OnRequest(func(req *colly.Request) {
		select {
		case <-ctx.Done():
			req.Abort()
		default:
		}
	})
	// Use custom Transport to cancel all pending requests at HTTP client,
	// which do not have a chance to stop at OnRequest callback.
	trans := &contextTransport{
		ctx:   ctx,
		trans: &http.Transport{},
	}
	c.WithTransport(trans)
}

type Result2 struct {
	Path         [][]string
	Degrees      int
	TimeTaken    time.Duration
	LinksVisited int
}

func IterativeDeepeningSearch(start, ends string) []byte {
	start = "https://en.wikipedia.org" + start
	ends = "https://en.wikipedia.org" + ends
	log.Println(start)
	log.Println(ends)
	targetUrl := ends

	c := colly.NewCollector(
		colly.Async(true),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() { cancel() }()

	collectorWithContext(c, ctx)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*wikipedia.*",
		Parallelism: 15,
		Delay:       100 * time.Millisecond,
	})

	// Set of paths to the target URL
	var results = set.NewSetOfSlice()
	maxDepth := 6
	depths := make(map[string]int)
	// Mutex to protect the depths map
	var depthsMutex sync.RWMutex
	// WaitGroup to wait for all URLs to be visited
	var wg sync.WaitGroup
	resultCh := make(chan *Result2)
	// Variable to store the least depth at which the target URL was found
	leastDepth := maxDepth + 1
	var result Result2
	// Counter for the number of links visited
	var linksVisited int64
	// Map to store the parent-child relationship of each URL
	parents := make(map[string]string)
	var targetFound bool
		// On every response
	c.OnRequest(func(r *colly.Request) {
		if result.Degrees != 0 {
			// Target found, cancel all further requests
			cancel()
		}
	})

	// Start the timer
	startTime := time.Now()

	// Channel to signal when the target URL is found
	done := make(chan struct{})

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Print the number of links visited every 10 seconds
	go func() {
		for range ticker.C {
			fmt.Println("Number of links visited:", atomic.LoadInt64(&linksVisited))
			fmt.Println("Time taken:", time.Since(startTime))
			
			if result.Degrees != 0 {
				cancel()
				return
			}
		}
	}()

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			cancel()
			return // Context canceled
		default:
			link := e.Attr("href")
			if isValidLink(link) {
				url := e.Request.AbsoluteURL(link)
				depthsMutex.RLock()
				currentDepth, ok := depths[e.Request.URL.String()]
				depthsMutex.RUnlock()
				if !ok {
					cancel()
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
							// Create a new instance of Result for this goroutine
							goroutineResult := Result2{
								Degrees:      currentDepth + 1,
								TimeTaken:    time.Since(startTime),
								LinksVisited: int(atomic.LoadInt64(&linksVisited)),
							}
							fmt.Println(url)
							fmt.Println("\nFound target URL at depth", leastDepth)
							fmt.Println("Time taken:", goroutineResult.TimeTaken)
							fmt.Println("Number of links visited:", goroutineResult.LinksVisited)
							path := []string{url}
							for url != start {
								parent, exists := parents[url]
								if !exists {
									break
								}
								path = append([]string{parent}, path...)
								url = parent
							}
							results.Add(path)
							goroutineResult.Path = results.ToSlice()
							// Send the result for this goroutine to the channel
							resultCh <- &goroutineResult
							targetFound = true
							close(done)

							cancel()
							return
						}
						if result.Degrees == 0 { // Target not found yet
							wg.Add(1)
							go func() {
								defer wg.Done()
								for {
									select {
									case <-ctx.Done():
										ctx.Done()
										return
									default:
										c.Visit(url)
									}
								}
							}()
							atomic.AddInt64(&linksVisited, 1)
						}
					} else {
						depthsMutex.Unlock()
					}
				}
			}
		}
        if targetFound {
            cancel() // Cancel the context if the target is found
        }
	})

	// Initialize the depths map with the starting URL
	depthsMutex.Lock()
	depths[start] = 0
	depthsMutex.Unlock()
	wg.Add(1)

	go func() {
		defer wg.Done()
		if ctx.Done() == nil {
			cancel()
			return
		}
		c.Visit(start)
	}()

	// Wait for the target URL to be found or all URLs to be visited
	select {
	case result := <-resultCh:
		var graph Result.Graph
		graph.GenerateGraph(results.ToSlice())

		outputResult := Result.Result{
			Traph:               graph,
			Time:                result.TimeTaken.Milliseconds(),
			TotalArticleChecked: result.LinksVisited,
			TotalArtcleVisit:    graph.GetNodesCount(),
		}

		jsonOutput, err := json.Marshal(outputResult)

		if err != nil {
			log.Println("Error marshalling JSON:", err)
			return []byte(fmt.Sprintf(`{Error marshalling JSON : %v}`, err))
		}

		log.Printf("JSON: %s\n", jsonOutput)
		return jsonOutput

	// Target URL found, return the result
	case <-ctx.Done():
		// Context canceled, return nil to indicate failure
		var graph Result.Graph
		graph.GenerateGraph(results.ToSlice())

		outputResult := Result.Result{
			Traph:               graph,
			Time:                result.TimeTaken.Milliseconds(),
			TotalArticleChecked: result.LinksVisited,
			TotalArtcleVisit:    graph.GetNodesCount(),
		}

		jsonOutput, err := json.Marshal(outputResult)

		if err != nil {
			log.Println("Error marshalling JSON:", err)
			cancel()
			return []byte(fmt.Sprintf(`{Error marshalling JSON : %v}`, err))
		}
		// Print json
		log.Printf("JSON: %s\n", jsonOutput)
		return jsonOutput
	}
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
