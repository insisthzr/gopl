package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/insisthzr/gopl/utils"
)

type LineCount struct {
	Count int
	Files utils.Set
}

func NewLineCount() LineCount {
	return LineCount{
		Files: utils.NewSet(),
	}
}

type LineCountsSync struct {
	lineCounts map[string]LineCount
	mu         sync.RWMutex
}

func NewLineCountsSync() LineCountsSync {
	return LineCountsSync{
		lineCounts: map[string]LineCount{},
	}
}

func (p *LineCountsSync) Counts(filename string, data string) {
	for _, line := range strings.Split(data, "\n") {
		p.Add(line, filename)
	}
}

func (p *LineCountsSync) Add(line string, file string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	_, ok := p.lineCounts[line]
	if !ok {
		p.lineCounts[line] = NewLineCount()
	}

	lineCount := p.lineCounts[line]
	newLineCount := LineCount{lineCount.Count + 1, lineCount.Files}
	p.lineCounts[line] = newLineCount
	p.lineCounts[line].Files.Add(file)
}

func (p *LineCountsSync) Get(line string) LineCount {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.lineCounts[line]
}

func (p *LineCountsSync) For(handler func(line string, lineCount LineCount)) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	//block write
	for line, lineCount := range p.lineCounts {
		handler(line, lineCount)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Printf("params too short")
	}

	ws := &sync.WaitGroup{}
	lineCounts := NewLineCountsSync()
	filenames := os.Args[1:]

	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			panic(err.Error())
		}
		defer f.Close()
		data, _ := ioutil.ReadAll(f)
		ws.Add(1)
		go func() {
			defer ws.Done()

			lineCounts.Counts(f.Name(), string(data))
		}()
	}
	ws.Wait()

	lineCounts.For(func(line string, lineCount LineCount) {
		if lineCount.Count > 1 {
			log.Printf("%d\t%s\t%s\n", lineCount.Count, line, lineCount.Files)
		}
	})
}
