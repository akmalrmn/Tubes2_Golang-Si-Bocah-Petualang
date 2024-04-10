package scraper

import (
	"github.com/hashicorp/golang-lru"
	"golang.org/x/net/context"
	"golang.org/x/net/html"
	"net"
	"net/http"
	"strings"
	"time"
)
var httpClient *http.Client
var linkCache, _ = lru.New(1000) // Cache for the links
func init() {
	httpClient = createCustomHTTPClient()
}

func createCustomHTTPClient() *http.Client {
	transport := &http.Transport{
		MaxIdleConns:       100,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	return &http.Client{
		Transport: transport,
	}
}

// ExtractLinks mengambil semua link dari halaman web Wikipedia
func ExtractLinks(url string) ([]string, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var links []string
	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			// End of the document, we're done
			return links, nil
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == "a" {
				for _, a := range t.Attr {
					if a.Key == "href" {
						trimmed := strings.Trim(a.Val, "\n")
						if !strings.HasPrefix(trimmed, "#") && strings.HasPrefix(trimmed, "/wiki/") && !strings.Contains(trimmed, ":") && !strings.Contains(trimmed, "Main_Page") && !strings.Contains(trimmed, "#") {
							links = append(links, trimmed)
						}
					}
				}
			}
		}
	}
}

// Fungsi rekursif untuk mengekstrak link dari node HTML

func ExtractLinksAsync(ctx context.Context, url string) ([]string, error) {
	// Check the cache before making an HTTP request
	if links, ok := linkCache.Get(url); ok {
		return links.([]string), nil
	}

	result := make(chan []string, 1)
	errResult := make(chan error, 1)

	go func() {
		defer close(result)
		defer close(errResult)
		links, err := ExtractLinks(url)
		if err != nil {
			result <- nil
			errResult <- err
			return
		}
		// Store the links in the cache
		linkCache.Add(url, links)
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
