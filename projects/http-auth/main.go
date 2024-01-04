package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
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

	http.HandleFunc("/authenticated", func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(os.Getenv("AUTH_USERNAME")))
			expectedPasswordHash := sha256.Sum256([]byte(os.Getenv("AUTH_PASSWORD")))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				w.Header().Add("Content-Type", "text/html")
				w.Write([]byte("<!DOCTYPE html><html>Hello username!</html>"))
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="localhost", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})

	http.ListenAndServe(":8080", nil)

}
