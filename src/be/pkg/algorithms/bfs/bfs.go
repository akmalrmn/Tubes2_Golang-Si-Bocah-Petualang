package bfs

import (
	"be/pkg/config"
	"be/pkg/scraper"
)

func BFS(starts, ends string, con *config.Config) []string {
	// Create a new Scraper
	s := scraper.NewScraper()
	s.SetScrapper(con)

	// Scrape the start URL
	urls := s.BFSScrape(starts, ends, con)

	return urls
}
