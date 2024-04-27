package scraper

import (
	"be/pkg/config"
	"be/pkg/set"
	"crypto/tls"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/proxy"
	"github.com/gocolly/colly/v2/queue"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	ArticleCount int = 0
)

func QueueColly(input, output *queue.Queue, start, ends string, con *config.Config, parents *sync.Map) *set.MapString {
	var results = set.NewSetOfSlice()

	// Instantiate default collector
	c := colly.NewCollector(
		colly.Async(con.IsAsync),
		colly.MaxDepth(con.MaxDepth),
		colly.UserAgent("Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_4_1) AppleWebKit/536.25 (KHTML, like Gecko) Chrome/51.0.2823.231 Safari/537"),
		colly.CacheDir(con.CacheDir),
		colly.AllowedDomains(con.AllowedDomains...),
	)

	extensions.RandomUserAgent(c)

	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // Timeout
			KeepAlive: 30 * time.Second, // keepAlive timeout
		}).DialContext,
		MaxIdleConns:          100,              // Maximum number of idle connections
		IdleConnTimeout:       90 * time.Second, // Idle connection timeout
		TLSHandshakeTimeout:   10 * time.Second, // TLS handshake timeout
		ExpectContinueTimeout: 1 * time.Second,
	})

	if p, err := proxy.RoundRobinProxySwitcher(
		"http://103.59.45.53:8080",
		"http://36.64.217.27:1313",
		"http://101.255.166.242:8080",
		"http://8c5f3b01ca4c7c3cd1c7ee83e4b652f69d6bd21b:@proxy.zenrows.com:8001"
	); err == nil {
		c.SetProxyFunc(p)
	}

	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: con.MaxParallelism,
		RandomDelay: con.RandomDelay,
	})

	if err != nil {
		log.Println("Error setting the limit rule:", err)
		return results
	}

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		url := e.Request.AbsoluteURL(link)

		if isValidLink(link) {
			// Check if link is visited
			_, visited := parents.Load(url)
			if visited {
				return
			}

			parentUrl := e.Request.URL.String()
			if url == parentUrl {
				return
			}
			parents.Store(url, parentUrl)

			// Add link to the second queue
			err = output.AddURL(url)
			if err != nil {
				log.Println("Error adding URL to queue:", err)
			}

		}

		// Check if the link is the destination
		if link == ends {

			path := []string{url}
			startFullPath := "https://en.wikipedia.org" + start

			for url != startFullPath {

				urlInterface, ok := parents.Load(url)
				if !ok {
					log.Println("Error: Link not found in parents map", url)
					return
				}
				url, ok = urlInterface.(string)
				if !ok {
					log.Println("Error: links is not a string")
					return
				}
				path = append([]string{url}, path...)
			}

			results.Add(path)
			return
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
	})

	c.OnRequest(func(r *colly.Request) {
		ArticleCount++
	})

	// Consume URLs in the first queue
	err = input.Run(c)
	if err != nil {
		log.Println("Error running the collector:", err)
		return results
	}

	c.Wait()
	return results
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
	return strings.HasPrefix(link, "/wiki/")
}
