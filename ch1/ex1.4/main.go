package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/insisthzr/gopl/utils"
	"io/ioutil"
	"strings"
	//"time"
)

type dup struct {
	count int
	files *utils.Set
}

func newDup() *dup {
	return &dup{files: utils.NewSet()}
}

type safeCounts struct {
	cs   map[string]*dup
	lock *sync.Mutex
}

func newSafeCounts() *safeCounts {
	return &safeCounts{cs: map[string]*dup{}, lock: new(sync.Mutex)}
}

func (p *safeCounts) counts(filename string, data string) {
	for _, line := range strings.Split(data, "\n") {
		p.add(line, filename)
	}
}

func (p *safeCounts) add(line string, file string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.cs[line] == nil {
		p.cs[line] = newDup()
	}
	p.cs[line].count++
	p.cs[line].files.Add(file)
}

func main() {
	ws := &sync.WaitGroup{}
	sc := newSafeCounts()
	files := os.Args[1:]

	for _, arg := range files {
		f, err := os.Open(arg)
		if err != nil {
			panic(err.Error())
		}
		defer f.Close()
		data, _ := ioutil.ReadAll(f)
		ws.Add(1)
		go func() {
			defer ws.Done()

			sc.counts(f.Name(), string(data))
		}()
	}

	ws.Wait()
	//time.Sleep(1 * time.Second)
	for line, d := range sc.cs {
		if d.count > 1 {
			fmt.Printf("%d\t%s\t%s\n", d.count, line, d.files)
		}
	}
}
