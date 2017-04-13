package main

import (
	"bufio"
	"fmt"
)

type WordCounter struct {
	words int
	lines int
}

func (c *WordCounter) Write(p []byte) (int, error) {
	words, err := c.countWords(p)
	if err != nil {
		return 0, err
	}
	c.words = words
	lines, err := c.countLines(p)
	if err != nil {
		return 0, err
	}
	c.lines = lines
	return len(p), nil
}

func (c *WordCounter) countWords(p []byte) (int, error) {
	count := 0
	for {
		next, word, err := bufio.ScanWords(p, true)
		if err != nil {
			return 0, err
		}
		if word == nil {
			break
		}
		p = p[next:]
		count++
	}
	return count, nil
}

func (c *WordCounter) countLines(p []byte) (int, error) {
	count := 0
	for {
		next, word, err := bufio.ScanLines(p, true)
		if err != nil {
			return 0, err
		}
		if word == nil {
			break
		}
		p = p[next:]
		count++
	}
	return count, nil
}

func main() {
	str := "hello world\nhello"
	c := &WordCounter{}
	fmt.Fprintf(c, str)
	fmt.Println(c)
}
