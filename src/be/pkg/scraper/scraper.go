package scraper

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

<<<<<<< Updated upstream
// ExtractLinks mengambil semua link dari halaman web Wikipedia
func ExtractLinks(url string) ([]string, error) {
	// Mendapatkan isi halaman web
	resp, err := http.Get(url)
=======
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
		"http://8c5f3b01ca4c7c3cd1c7ee83e4b652f69d6bd21b:@proxy.zenrows.com:8001",
	); err == nil {
		c.SetProxyFunc(p)
	}

	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: con.MaxParallelism,
		RandomDelay: con.RandomDelay,
	})

>>>>>>> Stashed changes
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
