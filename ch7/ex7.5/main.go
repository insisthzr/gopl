package main

import (
	"fmt"
	"io"
	"os"
)

type LimitedReader struct {
	R io.Reader
	N int64
}

func (r *LimitedReader) Read(p []byte) (n int, err error) {
	if r.N <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > r.N {
		p = p[:r.N]
	}
	n, err = r.R.Read(p)
	r.N -= int64(n)
	return
}

func LimitReader(r io.Reader, n int64) io.Reader {
	return &LimitedReader{r, n}
}

func main() {
	r := LimitReader(os.Stdin, 10)
	p := make([]byte, 5)
	for {
		_, err := r.Read(p)
		if err == io.EOF {
			break
		}
		fmt.Printf(string(p))
	}
	fmt.Printf("\n")
}
