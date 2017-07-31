package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	xkcdURL    = "https://xkcd.com/info.0.json"
	xkcdURLFmt = "https://xkcd.com/%d/info.0.json"

	usage    = "xkcd get N\nxkcd index\nxkcd search QUERY\n"
	filename = "index.json"
)

var (
	workerNum = 32
)

type Comic struct {
	Num        int    `json:"num"`
	Year       string `json:"year"`
	Month      string `json:"month"`
	Day        string `json:"day"`
	Title      string `json:"title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
}

type NumComicMap map[int]*Comic
type WordNumFoundMap map[string]map[int]bool

func getComicCount() (int, error) {
	resp, err := http.Get(xkcdURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status code %d", resp.StatusCode)
	}
	comic := &Comic{}
	err = json.NewDecoder(resp.Body).Decode(comic)
	if err != nil {
		return 0, err
	}
	return comic.Num, nil
}

func getComic(n int) (*Comic, error) {
	comic := &Comic{}
	url := fmt.Sprintf(xkcdURLFmt, n)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return comic, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return comic, fmt.Errorf("status code %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&comic)
	if err != nil {
		return comic, err
	}
	return comic, nil
}

func getComics() ([]*Comic, error) {
	max, err := getComicCount()
	if err != nil {
		return nil, err
	}
	comicChan := make(chan *Comic, max)
	numChan := make(chan int, max)

	go func() {
		mapNum(numChan, max) //map
		defer close(numChan)
	}()

	wg := &sync.WaitGroup{}
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go func() {
			consume(numChan, comicChan) //work
			defer wg.Done()
		}()
	}
	wg.Wait()
	close(comicChan)

	done := make(chan []*Comic)
	go reduce(comicChan, done) //reduce
	comics := <-done

	return comics, nil
}

func mapNum(numChan chan<- int, max int) {
	for i := 1; i <= max; i++ {
		numChan <- i
	}
}

func consume(numChan <-chan int, comicChan chan<- *Comic) {
	for n := range numChan {
		comic, err := getComic(n)
		if err != nil {
			fmt.Printf("getComic num %d error %v\n", n, err)
			continue
		}
		comicChan <- comic
	}
}

func reduce(comicChan <-chan *Comic, done chan<- []*Comic) {
	comics := []*Comic{}
	for comic := range comicChan {
		comics = append(comics, comic)
	}
	done <- comics
}

func genNumComicMap(comics []*Comic) NumComicMap {
	m := NumComicMap{}
	for _, comic := range comics {
		m[comic.Num] = comic
	}
	return m
}

func genWordNumFoundMap(comics []*Comic) WordNumFoundMap {
	m := WordNumFoundMap{}
	for _, comic := range comics {
		scanner := bufio.NewScanner(strings.NewReader(comic.Transcript))
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			word := strings.ToLower(scanner.Text())
			_, ok := m[word]
			if !ok {
				m[word] = map[int]bool{}
			}
			m[word][comic.Num] = true
		}
		err := scanner.Err()
		if err != nil {
			panic(err) //unexpected error
		}
	}
	return m
}

func index(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return writeIndex(file)
}

func writeIndex(writer io.Writer) error {
	comics, err := getComics()
	if err != nil {
		return err
	}
	numToComic := genNumComicMap(comics)
	wordToNum := genWordNumFoundMap(comics)

	enc := json.NewEncoder(writer)
	err = enc.Encode(numToComic)
	if err != nil {
		return err
	}
	err = enc.Encode(wordToNum)
	if err != nil {
		return err
	}
	return nil
}

func readIndex(reader io.Reader) (NumComicMap, WordNumFoundMap, error) {
	dec := json.NewDecoder(reader)
	numToComic := NumComicMap{}
	err := dec.Decode(&numToComic)
	if err != nil {
		return nil, nil, err
	}
	wordToNum := WordNumFoundMap{}
	err = dec.Decode(&wordToNum)
	if err != nil {
		return nil, nil, err
	}
	return numToComic, wordToNum, nil
}

func search(filename string, query string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	numToComic, wordToNum, err := readIndex(file)
	if err != nil {
		return err
	}
	getComicsByWords(numToComic, wordToNum, strings.Fields(query))
	return nil
}

func getComicsByWords(numToComic NumComicMap, wordToNum WordNumFoundMap, query []string) {
	found := map[int]int{}
	for _, word := range query {
		for num := range wordToNum[word] {
			found[num]++
		}
	}

	for num, count := range found {
		if count == len(query) { //match all words
			comic := numToComic[num]
			fmt.Printf("num %d total %s\n", num, comic.Title)
		}
	}
}

func printUsageAndExit() {
	fmt.Printf(usage)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		printUsageAndExit()
	}
	cmd := os.Args[1]
	switch cmd {
	case "get":
		if len(os.Args) != 3 {
			printUsageAndExit()
		}
		n, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("N (%s) must be an int\n", os.Args[1])
		}
		comic, err := getComic(n)
		if err != nil {
			fmt.Printf("Error getComic %v\n", err)
		}
		data, err := json.MarshalIndent(comic, " ", " ")
		if err != nil {
			fmt.Printf("Error getComic %v\n", err)
		}
		fmt.Println(string(data))
	case "index":
		if len(os.Args) != 2 {
			printUsageAndExit()
		}
		err := index(filename)
		if err != nil {
			fmt.Printf("Error index %v\n", err)
		}
	case "search":
		if len(os.Args) != 3 {
			printUsageAndExit()
		}
		query := os.Args[2]
		err := search(filename, query)
		if err != nil {
			fmt.Printf("Error index %v\n", err)
		}
	default:
		printUsageAndExit()
	}
}
