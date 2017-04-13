package main

import (
	"fmt"
	"io"
	"os"
)

type somethingCounter struct {
	w     io.Writer
	count int64
}

func (c *somethingCounter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	c.count += int64(n)
	return n, err
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	c := &somethingCounter{w, 0}
	return c, &c.count
}

func main() {
	writer, count := CountingWriter(os.Stdout)
	fmt.Fprintf(writer, "hello world\n")
	fmt.Println(*count)
	fmt.Fprintf(writer, "hello world\n")
	fmt.Println(*count)
}
