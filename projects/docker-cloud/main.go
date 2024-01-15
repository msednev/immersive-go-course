package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func run(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World."))
	})
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	log.Printf("Listening on %d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8090, "port the server listens to")
	flag.Parse()
	if err := run(port); err != nil {
		log.Fatal(err.Error())
	}

}