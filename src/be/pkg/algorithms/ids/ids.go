package ids

import (
    "context"
    "fmt"
    "log"
    "strings"
    "sync"
    "time"
    "net/http"
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
    // We can stop all requests at `OnRequest` callback 
    // before sending request to HTTP client.
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

type Result struct {
    Path          []string
    Degrees       int
    TimeTaken     time.Duration
    LinksVisited  int
}

func IterativeDeepeningSearch(start, ends string) *Result {
    targetUrl := ends

    c := colly.NewCollector(
        colly.Async(true),
    )

    ctx, cancel := context.WithCancel(context.Background())

    collectorWithContext(c, ctx)

    c.Limit(&colly.LimitRule{
        DomainGlob:  "*wikipedia.*",
        Parallelism: 15,
        Delay:       100 * time.Millisecond,
    })

    // The maximum depth for the IDS
    maxDepth := 6
    depths := make(map[string]int)
    // Mutex to protect the depths map
    var depthsMutex sync.RWMutex
    // WaitGroup to wait for all URLs to be visited
    var wg sync.WaitGroup
    resultCh := make(chan *Result)
    // Variable to store the least depth at which the target URL was found
    leastDepth := maxDepth + 1
    var result Result
    // Counter for the number of links visited
    var linksVisited int

    // Map to store the parent-child relationship of each URL
    parents := make(map[string]string)

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

    go func() {
        for range ticker.C {
            fmt.Println("Number of links visited:", linksVisited)
            fmt.Println("Time taken:", time.Since(startTime))
        }
    }()

    // On every a element which has href attribute call callback
    c.OnHTML("a[href]", func(e *colly.HTMLElement) {
        select {
        case <-ctx.Done():
            return // Context canceled
        default:
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
                            // Create a new instance of Result for this goroutine
                            goroutineResult := Result{
                                Degrees:      currentDepth + 1,
                                TimeTaken:    time.Since(startTime),
                                LinksVisited: linksVisited,
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
                            goroutineResult.Path = path
                            // Send the result for this goroutine to the channel
                            resultCh <- &goroutineResult
                            close(done)
                            return
                        }
                        if result.Degrees == 0 { // Target not found yet
                            wg.Add(1)
                            go func() {
                                defer wg.Done()
                                err := c.Visit(url)
                                if err != nil {
                                    log.Println("Error visiting URL:", url, "Error:", err)
                                }
                            }()
                            linksVisited++
                        }
                    } else {
                        depthsMutex.Unlock()
                    }
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

    // Wait for the target URL to be found or all URLs to be visited
    select {
    case result := <-resultCh:
        log.Println("Aaa")
        // Target URL found, return the result
        return result
    case <-ctx.Done():
        log.Println("Bbb")
        // Context canceled, return nil to indicate failure
        return &result}

    // Close the context to ensure all goroutines are stopped
    cancel()

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
