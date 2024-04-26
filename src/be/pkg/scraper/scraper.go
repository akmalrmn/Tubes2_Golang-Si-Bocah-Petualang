package scraper

import (
	"be/pkg/config"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
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
	s.colly = scrape.colly

	// Add the start URL to the queue
	s.queue.AddURL(starts)

	visited := &sync.Map{}

	// Process the queue
	for {
		// Process the queue
		done, q := s.processQueue(visited, ends)
		s.queue = q
		if done {
			break
		}
	}
	// Convert the visited sync.Map to a slice of strings
	var urls []string
	visited.Range(func(key, value interface{}) bool {
		urls = append(urls, key.(string))
		return true
	})

	return urls
}

func (s *Scraper) processQueue(visited *sync.Map, ends string) (bool, *queue.Queue) {

	s.colly.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// Process the queue
	s.colly.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		// Add the URL to the queue
		s.queue.AddURL(link)
		// Add the URL to the visited list
		visited.Store(link, struct{}{})
	})

	s.queue.Run(s.colly)
	s.colly.Wait()

	// Check if the target URL is found
	if _, ok := visited.Load(ends); ok {
		return true, s.queue
	}

	return false, s.queue

}
