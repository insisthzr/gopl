package main

import (
	"flag"
	"io"
	"log"
	"net"
	"time"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		_, err := io.WriteString(conn, time.Now().Format("15:04:05\n"))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

var (
	port = flag.String("port", "8000", "network port")
)

func main() {
	flag.Parse()
	addr := "localhost:" + *port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("started in %s\n", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}
