package scraper

import (
	"fmt"
	"github.com/hashicorp/golang-lru"
	"github.com/temoto/robotstxt"
	"golang.org/x/net/html"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	httpClient             *http.Client
	LinkCache, _           = lru.New(1000)
	ua                     *robotstxt.Group
	NumOfArticlesProcessed int
)

// Init inisialisasi HTTP client dan robots.txt
func Init() {
	httpClient = createCustomHTTPClient()
	resp, err := httpClient.Get("https://en.wikipedia.org/robots.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	data, err := robotstxt.FromResponse(resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	ua = data.FindGroup("*")
	if ua == nil {
		fmt.Println("No group found for user agent")
		return
	}
}

// createCustomHTTPClient membuat HTTP client dengan konfigurasi khusus
func createCustomHTTPClient() *http.Client {
	transport := &http.Transport{
		MaxIdleConns:       300,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}).DialContext,
	}

	return &http.Client{
		Transport: transport,
	}
}

// ExtractLinks mengambil semua link dari halaman web Wikipedia
func ExtractLinks(url string) ([]string, error) {

	if !ua.Test(url) {
		return nil, fmt.Errorf("URL is not allowed by robots.txt")
	}

	if links, ok := LinkCache.Get(url); ok {
		return links.([]string), nil
	} // Check the cache before making an HTTP request

	resp, err := httpClient.Get(url)
	NumOfArticlesProcessed++
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	var links []string
	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			LinkCache.Add(url, links)
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
		default:
			continue
		}
	}
}

// ExtractLinksNonAdd mengambil semua link dari halaman web Wikipedia tanpa menambahkan counter artikel yang telah di proses
func ExtractLinksNonAdd(url string) ([]string, error) {
	if links, ok := LinkCache.Get(url); ok {
		return links.([]string), nil
	}
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	var links []string
	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			LinkCache.Add(url, links)
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
		default:
			continue
		}
	}
}

// ExtractLinksAsync mengambil semua link dari halaman web Wikipedia secara asynchronous
func ExtractLinksAsync(url string) ([]string, error) {

	// Check the cache before making an HTTP request
	if links, ok := LinkCache.Get(url); ok {
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
		LinkCache.Add(url, links)
		result <- links
	}()

	select {
	case links := <-result:
		return links, nil
	case err := <-errResult:
		return nil, err
	}
}

func linkToTitle(link string) string {
	link = strings.TrimPrefix(link, "/wiki/")
	parts := strings.Split(link, "/")
	replaced := strings.Replace(parts[len(parts)-1], "_", " ", -1)
	return replaced
}

func titleToLink(title string) string {
	replaced := strings.Replace(title, " ", "_", -1)
	return "/wiki/" + replaced
}
