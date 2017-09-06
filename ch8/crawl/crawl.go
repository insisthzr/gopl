package main

import (
	"flag"
	"log"

	"github.com/insisthzr/gopl/ch5/links"
)

var (
	tokens = make(chan struct{}, 20)
)

type data struct {
	links []string
	depth int
}

func crawl(url string) []string {
	log.Println(url)
	tokens <- struct{}{}
	defer func() {
		<-tokens
	}()
	list, err := links.Extract(url)
	if err != nil {
		log.Println(err)
	}
	return list
}

var (
	depth = flag.Int("d", 3, "search depth")
	url   = flag.String("u", "http://gopl.io", "url")
)

func main() {
	flag.Parse()

	worklist := make(chan data, 1)
	worklist <- data{links: []string{*url}, depth: 0}
	n := 0
	n++

	seen := make(map[string]bool)
	for ; n > 0; n-- {
		d := <-worklist
		if d.depth >= *depth {
			continue
		}
		for _, link := range d.links {
			if !seen[link] {
				seen[link] = true
				n++ // a mark for I'm reading
				go func(link string) {
					links := crawl(link)
					worklist <- data{links: links, depth: d.depth + 1}
				}(link)
			}
		}
	}
}
