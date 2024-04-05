package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"strings"
	"sync"
)

var linkCache sync.Map // Cache for the links

// ExtractLinks mengambil semua link dari halaman web Wikipedia
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

// Fungsi rekursif untuk mengekstrak link dari node HTML

func ExtractLinksAsync(ctx context.Context, url string) ([]string, error) {
	// Check the cache before making an HTTP request
	if links, ok := linkCache.Load(url); ok {
		return links.([]string), nil
	}

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
		// Store the links in the cache
		linkCache.Store(url, links)
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

func makeUnique(links []string) []string {
	keys := make(map[string]bool)
	var uniqueLinks []string
	for _, link := range links {
		if _, value := keys[link]; !value {
			keys[link] = true
			uniqueLinks = append(uniqueLinks, link)
		}
	}
	return uniqueLinks
}

func linkToTitle(link string) string {
	parts := strings.Split(link, "/")
	replaced := strings.Replace(parts[len(parts)-1], "_", " ", -1)

	return replaced
}

func titleToLink(title string) string {
	replaced := strings.Replace(title, " ", "_", -1)

	return "/wiki/" + replaced
}
