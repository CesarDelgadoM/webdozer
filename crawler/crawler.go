package crawler

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/CesarDelgadoM/webdozer/database"
	"github.com/CesarDelgadoM/webdozer/utils"
	"golang.org/x/net/html"
)

const (
	label     string = "a"
	attribute string = "href"
	depth     int    = 3
	key       string = "cache-"
)

type Crawler struct {
	urlBase   string
	cache     database.Repository
	wg        *sync.WaitGroup
	results   chan *Result
	UrlsFound chan string
}

type Result struct {
	url   string
	urls  []string
	depth int
}

func NewCrawler(wg *sync.WaitGroup, redis *database.RedisPool) *Crawler {

	return &Crawler{
		cache:     database.NewSet(redis),
		wg:        wg,
		results:   make(chan *Result, 10000),
		UrlsFound: make(chan string, 10000),
	}
}

func (c *Crawler) LaunchCrawler(url string) {

	c.urlBase = url
	cacheKey := key + utils.ExtractNameUrl(url)

	go c.Crawl(url, depth)
	go func() {
		defer c.wg.Done()
		var m sync.Mutex

		for {
			select {
			case r := <-c.results:

				c.UrlsFound <- r.url

				for _, url := range r.urls {
					m.Lock()
					if !c.cache.Exist(cacheKey, url) && r.depth > 0 {
						go c.Crawl(url, r.depth-1)
						err := c.cache.Add(cacheKey, url)
						if err != nil {
							log.Println("Error caching url:", err)
						}
					}
					m.Unlock()
				}

			case <-time.After(10 * time.Second):
				if len(c.results) == 0 {
					log.Println("Channel results is empty, waiting 30 seconds...")
					<-time.After(30 * time.Second)

					if len(c.results) == 0 {
						c.cache.Del(cacheKey)
						log.Println("Deleted cache:", cacheKey)
						close(c.results)
						close(c.UrlsFound)
						return
					}
				}
			}
		}
	}()
}

func (c *Crawler) Crawl(url string, depth int) {

	urls, err := c.getUrlsFromPage(url)
	if err != nil {
		url = c.urlBase + url
		urls, err = c.getUrlsFromPage(url)
		if err != nil {
			return
		}
	}

	c.results <- &Result{
		url:   url,
		urls:  urls,
		depth: depth,
	}
}

func (c *Crawler) getUrlsFromPage(url string) ([]string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed parsing %s as HTML: %v", url, err)
	}
	defer resp.Body.Close()

	var urls []string = make([]string, 0)
	return c.extractUrls(urls, doc), nil
}

func (c *Crawler) extractUrls(urls []string, doc *html.Node) []string {

	if doc.Type == html.ElementNode && doc.Data == label {
		for _, a := range doc.Attr {
			if a.Key == attribute {
				urls = append(urls, a.Val)
			}
		}
	}

	for n := doc.FirstChild; n != nil; n = n.NextSibling {
		urls = c.extractUrls(urls, n)
	}

	return urls
}
