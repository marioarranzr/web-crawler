package crawler

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/marioarranzr/web-crawler/parser"
)

type SiteMap struct {
	Pages map[string]parser.PageDetails `json:"pages"`
}

type namedPageDetails struct {
	url     *url.URL
	details parser.PageDetails
}

type Crawler struct {
	startURL *url.URL
	siteMap  *SiteMap
	idx      index
	concur   int
}

func New(u *url.URL, concur int) *Crawler {
	siteMap := &SiteMap{Pages: make(map[string]parser.PageDetails)}
	idx := make(index)
	idx.add(u)
	c := &Crawler{u, siteMap, idx, concur}
	return c
}

func (c *Crawler) Crawl() *SiteMap {
	count := 0
	pagesToVisit := make(chan *url.URL, c.concur)
	results := make(chan namedPageDetails, c.concur)
	for i := 0; i < c.concur; i++ {
		go crawlPage(c.startURL, pagesToVisit, results)
	}
	for {
		select {
		case result := <-results:
			for _, link := range result.details.InternalLinks {
				linkURL, _ := url.Parse(link)
				c.idx.add(linkURL)
			}
			c.siteMap.Pages[result.url.String()] = result.details
			count--
		default:
			unvisitedLinks := c.idx.getUnvisitedLinks()
			numUnvisitedLinks := len(unvisitedLinks)
			if numUnvisitedLinks == 0 && count == 0 {
				return c.siteMap
			} else if numUnvisitedLinks > 0 {
				l := unvisitedLinks[0]
				c.idx.markVisited(l)
				pagesToVisit <- l
				count++
			}
		}
	}
}

func crawlPage(startURL *url.URL, urls <-chan *url.URL, results chan<- namedPageDetails) {
	for {
		select {
		case u := <-urls:
			body, err := getBody(u)
			if err != nil {
				log.Printf("Error reading webpage '%s': %v", u, err)
				return
			}
			details := parser.Parse(startURL, bytes.NewReader(body))
			namedDetails := namedPageDetails{
				url:     u,
				details: details,
			}
			results <- namedDetails
		default:
		}
	}
}

type index map[string]bool

func (i index) add(u *url.URL) {
	u.Fragment = ""
	u.RawQuery = ""
	if _, ok := i[u.String()]; !ok {
		log.Printf("Adding %v to index", u)
		i[u.String()] = false
	}
}

func (i index) markVisited(u *url.URL) {
	u.Fragment = ""
	u.RawQuery = ""
	log.Printf("Marking %v as visited", u)
	i[u.String()] = true
}

func (i index) getUnvisitedLinks() []*url.URL {
	unvisitedLinks := []*url.URL{}
	for k, v := range i {
		if !v {
			u, _ := url.Parse(k)
			unvisitedLinks = append(unvisitedLinks, u)
		}
	}
	return unvisitedLinks
}

func getBody(u *url.URL) ([]byte, error) {
	log.Printf("Fetching HTML from '%s'", u)
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
