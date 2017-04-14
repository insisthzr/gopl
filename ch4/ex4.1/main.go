package main

import (
	"fmt"
	"golang.org/x/net/html"
	"os"
	//"io"
	"io"
	//"io/ioutil"
	"net/http"
)

func main() {
	body := fetch(os.Args[1])
	defer body.Close()

	//result, _ := ioutil.ReadAll(body)
	//fmt.Println(string(result))

	doc, err := html.Parse(body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
		os.Exit(1)
	}
	for _, link := range visit(nil, doc) {
		fmt.Println(link)
	}
}

func fetch(url string) io.ReadCloser {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}
	return resp.Body
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				fmt.Println(a.Val)
				links = append(links, a.Val)
			}
		}
	}

	//for c := n.FirstChild; c != nil; c = c.NextSibling {
	//	links = visit(links, c)
	//}

	links = doVisit(links, n.FirstChild)

	return links
}

func doVisit(links []string, n *html.Node) []string {
	if n == nil {
		return links
	}

	links = visit(links, n)
	links = doVisit(links, n.NextSibling)
	return links
}
