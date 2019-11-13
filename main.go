package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/marioarranzr/web-crawler/crawler"
)

var concur = flag.Int("concurrency", 10, "Number of concurrent requests.")
var output = flag.String("output", "output", "output name file.")

func init() {
	usage()
}

func usage() {
	fmt.Fprintln(os.Stderr,
		"Usage: web-crawler [-concurrency n] [-output output] http://www.example.com")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	var (
		u   *url.URL
		f   *os.File
		c   *crawler.Crawler
		err error
	)
	flag.Parse()
	rawurl := flag.Arg(0)
	u, err = url.ParseRequestURI(rawurl)
	if err != nil {
		message := fmt.Sprintf(
			"Could not validate url '%s'.\n%v.\n", rawurl, err)
		usage()
		fmt.Fprintln(os.Stderr, "Error:", message)
		os.Exit(1)
	}
	c = crawler.New(u, *concur)
	siteMap := c.Crawl()
	var jsonSiteMap []byte
	jsonSiteMap, err = json.MarshalIndent(siteMap, "", "  ")
	if err != nil {
		usage()
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	f, err = os.Create(fmt.Sprintf("%s.txt", *output))
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.WriteString(string(jsonSiteMap))
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
