// Fetch prints the content found at a URL.
package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		panic("too few args")
	}

	url := os.Args[1]
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	resp, err := http.Get(url)
	if err != nil {
		panic(errors.Wrap(err, fmt.Sprintf("http.Get(%q)", url)))
	}
	defer resp.Body.Close()

	fmt.Fprintf(os.Stdout, "statuscode: %d\n", resp.StatusCode)
	io.Copy(os.Stdout, resp.Body)
}
