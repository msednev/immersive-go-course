package main

import (
	"context"
	"encoding/json"
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
	URL     string `json:"url"`
}

// var images = []Image{
// 	{
// 		Title: "Sunset",
// 		AltText: "Clouds at sunset",
// 		URL: "https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
// 	},
// 	{
// 		Title: "Mountain",
// 		AltText: "A mountain at sunset",
// 		URL: "https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
// 	},
// }

func main() {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		fmt.Fprintln(os.Stderr, "Provide a database URL")
		os.Exit(1)
	}

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable connect to database: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := conn.Close(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "unable to close database connection: %v", err)
			os.Exit(1)
		}
	}()

	rows, err := conn.Query(context.Background(), "SELECT title, url, alt_text FROM public.images")
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to execute query: %v", err)
		os.Exit(1)
	}

	var images []Image
	var title, url, altText string

	for rows.Next() {
		if err = rows.Scan(&title, &url, &altText); err != nil {
			fmt.Fprintf(os.Stderr, "unable to extract columns: %v", err)
			os.Exit(1)
		}
		images = append(images, Image{Title: title, URL: url, AltText: altText})
	}

	http.HandleFunc("/images.json", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		indent, err := strconv.Atoi(queryParams.Get("indent"))
		if err != nil {
			http.Error(w, "Indent should be an integer", http.StatusBadRequest)
		}
		b, err := json.MarshalIndent(images, "", strings.Repeat(" ", indent))
		if err != nil {
			http.Error(w, "Cannot serialize object", http.StatusInternalServerError)
		}

		w.Header().Add("Content-Type", "text/json")
		w.Write([]byte(b))
	})

	http.ListenAndServe(":8080", nil)
}
