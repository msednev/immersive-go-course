package main

import (
	"fmt"
	"html"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" && len(r.URL.Query()) != 0 {
			w.Header().Add("Content-Type", "text/html")
			w.Write([]byte("<!DOCTYPE html><html><em>Hello, world</em><p>Query parameters:<ul>"))
			for k, v := range r.URL.Query() {
				output := html.EscapeString(fmt.Sprintf("<li>%v: %v</li>", k, v))
				w.Write([]byte(output))
			}
			w.Write([]byte("</li></ul></p></html>"))

		} else if r.Method == "GET" {
			w.Header().Add("Content-Type", "text/html")
			w.Write([]byte("<!DOCTYPE html><html><em>Hello, world</em></html>"))

		} else if r.Method == "POST" {
			content, _ := io.ReadAll(r.Body)
			w.Write(content)
		}
	})

	http.HandleFunc("/200", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("200 OK"))
	})

	http.Handle("/404", http.NotFoundHandler())

	http.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("500 Internal Server Error"))
	})

	http.ListenAndServe(":8080", nil)

}
