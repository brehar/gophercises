package main

import (
	"../linkparser"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	urlFlag := flag.String("url", "http://gophercises.com", "The URL for which you would like to build a sitemap")
	maxDepth := flag.Int("depth", 3, "The maximum number of levels to traverse")

	flag.Parse()

	pages := bfs(*urlFlag, *maxDepth)
	toXml := urlset{
		Urls:  make([]loc, len(pages)),
		Xmlns: xmlns,
	}

	for i, page := range pages {
		toXml.Urls[i] = loc{page}
	}

	fmt.Print(xml.Header)

	enc := xml.NewEncoder(os.Stdout)

	enc.Indent("", "  ")

	if err := enc.Encode(toXml); err != nil {
		panic(err)
	}

	fmt.Println()
}

type empty struct{}

func bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]empty)

	var q map[string]empty

	nq := map[string]empty{
		urlStr: {},
	}

	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]empty)

		if len(q) == 0 {
			break
		}

		for n := range q {
			if _, ok := seen[n]; ok {
				continue
			}

			seen[n] = empty{}

			for _, link := range get(n) {
				if _, ok := seen[link]; !ok {
					nq[link] = empty{}
				}
			}
		}
	}

	ret := make([]string, 0, len(seen))

	for site := range seen {
		ret = append(ret, site)
	}

	return ret
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)

	if err != nil {
		return []string{}
	}

	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()

	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	links, _ := linkparser.Parse(r)

	var ret []string

	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}

	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string

	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}

	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}
