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
	ArticleCount = 0
)

func QueueColly(input, output *queue.Queue, start, ends string, con *config.Config, parents *sync.Map) *set.MapString {
	var results = set.NewSetOfSlice()

	// Instantiate default collector
	c := colly.NewCollector(
		colly.Async(con.IsAsync),
		colly.MaxDepth(con.MaxDepth),
		colly.UserAgent("Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-TW; rv:1.9.2.4) Gecko/20100611 Firefox/3.6.4 GTB7.0 ( .NET CLR 3.5.30729)"),
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
		IdleConnTimeout:       90 * time.Second, // Idle connection timeout
		TLSHandshakeTimeout:   10 * time.Second, // TLS handshake timeout
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	})

	// Round-robin Proxy
	rp, err := proxy.RoundRobinProxySwitcher(
		"http://scrapingant&proxy_type=residential&proxy_country=ID&browser=false:cabfe254df5746b4bdc0cdcd762fc034@proxy.scrapingant.com:8080",
		"http://117.54.114.101:80",
		"http://58.147.189.222:3128",
		"http://43.133.136.208:8800",
		"http://103.105.196.128:80",
		"http://117.54.114.99:80",
	)
	if err != nil {
		log.Println("Error setting the proxy:", err)
	}
	c.SetProxyFunc(rp)

	err = c.Limit(&colly.LimitRule{
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
		log.Println("Visiting", r.URL.String())
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
