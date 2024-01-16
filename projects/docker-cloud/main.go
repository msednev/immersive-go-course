package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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

func getPort() (int, error) {
	if os.Getenv("HTTP_PORT") == "" {
		return 80, nil
	}

	port, err := strconv.Atoi(os.Getenv("HTTP_PORT"))
	if err != nil {
		return 0, fmt.Errorf("provide a valid port: %w", err)
	}
	return port, nil
}

func main() {
	port, err := getPort()
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := run(port); err != nil {
		log.Fatal(err.Error())
	}
}