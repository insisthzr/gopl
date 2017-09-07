package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	timeout = 10 * time.Second
)

func echo(c net.Conn, shout string, delay time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println(shout)
	fmt.Fprintf(c, "%s\n", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintf(c, "%s\n", shout)
	time.Sleep(delay)
	fmt.Fprintf(c, "%s\n", strings.ToLower(shout))
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	readChan := make(chan string)
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	wg := &sync.WaitGroup{}
	go recv(conn, readChan)
LOOP:
	for {
		select {
		case text := <-readChan:
			timer.Reset(timeout)
			wg.Add(1)
			go echo(conn, text, 5*time.Second, wg)
		case <-timer.C:
			log.Println("timeout")
			break LOOP
		}
	}
	wg.Wait()
	log.Println("done")
}

func recv(conn net.Conn, in chan<- string) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)
		in <- text
	}
	err := scanner.Err()
	if err != nil {
		log.Println(err)
	}
}

func main() {
	listner, err := net.Listen("tcp", ":20000")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("listening at 20000")
	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("client connected")
		go handleConn(conn)
	}
}
