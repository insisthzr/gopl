package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

const (
	usage = "findlinks {url}"
)

var (
	errNot200 = errors.New("not 200")
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}
	url := os.Args[1]
	rc, err := fetch(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rc.Close()
	doc, err := html.Parse(rc)
	if err != nil {
		panic(err)
	}
	links := visist(nil, doc)
	for _, link := range links {
		fmt.Println(link)
	}
}

func visist(links []string, n *html.Node) []string {
	if n == nil {
		return links
	}

	if n.Type == html.ElementNode {
		if n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, a.Val)
				}
			}
		}
	}

	links = visist(links, n.FirstChild)
	links = visist(links, n.NextSibling)

	return links
}

func fetch(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errNot200
	}
	return resp.Body, nil
}
