package main

import (
	"log"
	"os"

	"github.com/insisthzr/gopl/ch5/links"
)

func bfs(fn func(item string) []string, worklist []string) {
	seen := map[string]bool{}
	for len(worklist) > 0 {
		items := worklist
		worklist = []string{}
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, fn(item)...)
			}
		}
	}
}

func crawl(url string) []string {
	log.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Println(err)
	}
	return list
}

func main() {
	bfs(crawl, os.Args[1:])
}
