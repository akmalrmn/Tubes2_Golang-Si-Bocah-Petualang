package scraper

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// ExtractLinks mengambil semua link dari halaman web Wikipedia
func ExtractLinks(url string) ([]string, error) {
	// Mendapatkan isi halaman web
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("gagal melakukan GET request: %v", err)
	}
	defer resp.Body.Close()

	// Parsing HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal melakukan parsing HTML: %v", err)
	}

	// Mengambil semua link dari dokumen
	links := extractLinks(doc)

	return links, nil
}

// Fungsi rekursif untuk mengekstrak link dari node HTML
func extractLinks(n *html.Node) []string {
	var links []string
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				link := strings.TrimSpace(attr.Val)
				
				// Pastikan link bukan merupakan link internal Wikipedia (misalnya, tidak dimulai dengan "#")
				if !strings.HasPrefix(link, "#") && strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":"){
					links = append(links, link)
				}
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, extractLinks(c)...)
	}
	return makeUnique(links)
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
