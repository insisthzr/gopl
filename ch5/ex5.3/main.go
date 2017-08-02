package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"

	"github.com/insisthzr/gopl/utils"
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
	err = tagText(rc, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func tagText(reader io.Reader, writer io.Writer) error {
	z := html.NewTokenizer(reader)
	stack := utils.NewStack()
LOOP:
	for {
		typ := z.Next()
		switch typ {
		case html.ErrorToken:
			break LOOP
		case html.StartTagToken:
			name, _ := z.TagName()
			stack.Push(string(name))
		case html.EndTagToken:
			_, ok := stack.Pop()
			if !ok {
				//return fmt.Errorf("no start tag before end tag: %v", name)
			}
		case html.TextToken:
			text := z.Text()
			if len(strings.TrimSpace(string(text))) == 0 {
				continue
			}
			cur, ok := stack.Pop()
			if !ok {
				continue
				//return fmt.Errorf("no start tag before text tag: %s", string(text))
			}
			name, ok := cur.(string)
			if !ok {
				panic("interface{} must be string")
			}
			if name == "script" || name == "style" {
				continue
			}
			writer.Write([]byte(fmt.Sprintf("<%s>", name)))
			writer.Write([]byte(text))
			writer.Write([]byte("\n"))
		}
	}
	return nil
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
