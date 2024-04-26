package scraper

import (
	"be/pkg/config"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"net/url"
	"sync"
)

// Scrape Make new interface for Scraper
type Scrape interface {
	BFSScrape(starts, ends string, con config.Config) []string
}

// Scraper Create a struct for Scraper
type Scraper struct {
	colly *colly.Collector
	queue *queue.Queue
}

// NewScraper Create a new Scraper
func NewScraper() *Scraper {
	q, _ := queue.New(2, &queue.InMemoryQueueStorage{MaxSize: 10000}) // handle the error
	return &Scraper{
		colly: colly.NewCollector(),
		queue: q,
	}
}

func (s *Scraper) SetScrapper(config *config.Config) {
	s.colly = colly.NewCollector(
		colly.Async(config.IsAsync),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
		colly.MaxDepth(config.MaxDepth),
		colly.AllowedDomains(config.AllowedDomains...),
	)

	s.colly.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: config.MaxParallelism,
	})
}

// BFSScrape Scrape the URL using BFS method
func (s *Scraper) BFSScrape(starts, ends string, con *config.Config) []string {
	// Create a new Collector
	scrape := NewScraper()
	scrape.SetScrapper(con)

	// Add the start URL to the queue
	scrape.queue.AddURL(starts)

	visited := sync.Map{}

	// Make colly request
	req := &colly.Request{URL: &url.URL{Path: starts}}

	// Process the queue

	found , OutQueue := scrape.processQueue(req, visited, ends)

	for !found {
		// Get the next URL from the queue
		s.queue = OutQueue
		// Process the queue
		found , OutQueue = scrape.processQueue(req, visited, ends)
	}

	// Return the path
	return []string{}
}
func (s *Scraper) processQueue(node *colly.Request, visited sync.Map, ends shu
tring) (bool, *queue.Queue) {
	// Mark the URL as visited
	fmt.Println("Visiting", node.URL.String())
	visited.Store(node.URL.String(), true)

	// Find the links on the page
	s.colly.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if link == ends {
			// Found the end URL
			// Print the path
			visited.Range(func(key, value interface{}) bool {
				println(key.(string))
			})
			return true , _
		}
		// Check if the link is already visited
		if _, ok := visited.Load(link); !ok {
			// Add the link to the queue
			nextURL, err := url.Parse(e.Request.AbsoluteURL(link))
			if err != nil {
				// handle error
				fmt.Println("Invalid URL:", err)
				return
			}
			nextRequest := &colly.Request{URL: nextURL}
			s.processQueue(nextRequest, visited, ends)
		}
	})

	// Visit the URL
	s.colly.Visit(node.URL.String())

	// Wait for the request to finish
	s.colly.Wait()

	return false, s.queue
}
