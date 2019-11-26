package parser

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"net/url"
)

type PageDetails struct {
	InternalLinks []string `json:"internal_links"`
	ExternalLinks []string `json:"external_links"`
	Assets        []string `json:"assets"`
}

func Parse(pageURL *url.URL, webpage io.Reader) PageDetails {
	details := PageDetails{}
	tokenizer := html.NewTokenizer(webpage)
	for {
		if tokenizer.Next() == html.ErrorToken {
			return details
		}
		token := tokenizer.Token()
		if token.Type != html.StartTagToken {
			continue
		}
		switch token.DataAtom {
		case atom.Link:
			rawurl := getHref(token.Attr)
			if len(rawurl) == 0 {
				continue
			}
			u, _ := url.Parse(rawurl)
			resolvedURL := pageURL.ResolveReference(u)
			details.Assets = append(details.Assets, resolvedURL.String())
		case atom.A:
			rawurl := getHref(token.Attr)
			if len(rawurl) == 0 {
				continue
			}
			u, _ := url.Parse(rawurl)
			resolvedURL := pageURL.ResolveReference(u)
			if resolvedURL.Host == pageURL.Host {
				details.InternalLinks = append(details.InternalLinks,
					resolvedURL.String())
			} else {
				details.ExternalLinks = append(details.ExternalLinks,
					resolvedURL.String())
			}
		}
	}
}

func getHref(a []html.Attribute) string {
	for _, attr := range a {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}
