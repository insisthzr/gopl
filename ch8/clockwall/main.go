package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

type Clock struct {
	name string
	host string
}

func (c *Clock) watch(reader io.Reader, writer io.Writer) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Fprintf(writer, "%s %s\n", c.name, scanner.Text())
	}
	err := scanner.Err()
	if err != nil {
		log.Printf("%v closed because %v\n", c, err)
	} else {
		log.Printf("%v closed\n", c)
	}
}

func parseHosts(str string) ([]*Clock, error) {
	hosts := strings.Split(str, ",")
	clocks := make([]*Clock, 0, len(hosts))
	for _, host := range hosts {
		nameHost := strings.Split(host, "=")
		if len(nameHost) != 2 {
			return nil, fmt.Errorf("invalid name host pair: %s", host)
		}
		clocks = append(clocks, &Clock{name: nameHost[0], host: nameHost[1]})
	}
	return clocks, nil
}

var (
	hosts = flag.String("hosts", "", "NAME=HOST[,NAME=HOST]")
)

func main() {
	flag.Parse()
	if hosts == nil {
		log.Fatalf("-hosts is required\n")
	}

	clocks, err := parseHosts(*hosts)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	wg := &sync.WaitGroup{}
	for _, clock := range clocks {
		wg.Add(1)
		go func(clock *Clock) {
			defer wg.Done()
			conn, err := net.Dial("tcp", clock.host)
			if err != nil {
				log.Println(err)
				return
			}
			clock.watch(conn, os.Stdout)
		}(clock)
	}
	wg.Wait()
	log.Printf("done\n")
}
