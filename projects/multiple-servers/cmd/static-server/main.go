package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"servers/static"
)

func main() {
	pathPtr := flag.String("path", "", "where to read static files from")
	portPtr := flag.Int("port", 0, "a port this server will be running on")
	flag.Parse()
	static.Run(*pathPtr, *portPtr)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html, err := os.Open("./assets/index.html")
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to open file: %v", err)
			os.Exit(1)
		}
		body, err := io.ReadAll(html)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to read content: %v", err)
			os.Exit(1)
		}
		if _, err := w.Write(body); err != nil {
			fmt.Fprintf(os.Stderr, "unable to write request body: %v", err)
			os.Exit(1)
		}
	})

	http.ListenAndServe(":8082", mux)
}
