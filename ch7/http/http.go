package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(404)
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
	fmt.Fprintf(w, "%f\n", price)
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if item == "" {
		w.WriteHeader(400)
		fmt.Fprintf(w, "item must be not nil\n")
		return
	}
	priceStr := req.URL.Query().Get("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}
	db[item] = dollars(price)
	fmt.Fprintf(w, "ok")
}

func main() {
	db := database{"shoes": 50, "socks": 5}
	mux := http.NewServeMux()
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)
	mux.HandleFunc("/update", db.update)
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}
