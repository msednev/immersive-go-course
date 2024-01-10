package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
)

type Image struct {
	Title   string `json:"title"`
	AltText string `json:"alt_text"`
	Url     string `json:"url"`
}


func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to create to database connection: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := conn.Close(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "unable to close database connection: %v", err)
			return
		}
	}()

	portPtr := flag.String("port", "", "a port the server listens to")
	flag.Parse()

	var images []Image
	var title, url, altText string
	
	mux := http.NewServeMux()

	mux.HandleFunc("/images.json", func(w http.ResponseWriter, r *http.Request) {
		indentPar := r.URL.Query().Get("indent")
		indent := ""
		if indentPar != "" {
			multiplier, err := strconv.Atoi(indentPar)
			if err != nil {
				fmt.Fprintf(os.Stderr, "indent should be integer: %v", err)
				http.Error(w, "provided a non-int value for indent", http.StatusBadRequest)
				return
			}
			indent = strings.Repeat(indent, multiplier)
		}

		if r.Method == "GET" {
			rows, err := conn.Query(context.Background(), "SELECT title, alt_text, url FROM images")
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to get images: %v", err)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			for rows.Next() {
				if err := rows.Scan(&title, &altText, &url); err != nil {
					fmt.Fprintf(os.Stderr, "unable to get column row values: %v", err)
					return
				}
				images = append(images, Image{Title: title, AltText: altText, Url: url})
			}

			payload, err := json.MarshalIndent(images, "", indent)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to marshall to json: %v", err)
				http.Error(w, "", http.StatusInternalServerError)
			}

			w.Header().Set("Content-Type", "text/json")
			if _, err := w.Write(payload); err != nil {
				fmt.Fprintf(os.Stderr, "unable to write response: %v", err)
			}
		}
	})

	if err := http.ListenAndServe(":" + *portPtr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "unable to start server: %v", err)
		os.Exit(1)
	}
}
