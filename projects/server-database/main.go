package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

func fetchImages(conn *pgx.Conn) ([]Image, error) {
	var images []Image
	var title, url, altText string

	rows, err := conn.Query(context.Background(), "SELECT title, url, alt_text FROM public.images")
	if err != nil {
		return nil, fmt.Errorf("unable to execute query: %v", err)
	}

	for rows.Next() {
		if err = rows.Scan(&title, &url, &altText); err != nil {
			return nil, fmt.Errorf("unable to extract columns: %v", err)
		}
		images = append(images, Image{Title: title, URL: url, AltText: altText})
	}

	return images, nil
	
}

func addImage(conn *pgx.Conn, image Image) error {
	_, err := conn.Exec(
		context.Background(),
		"INSERT INTO public.images (title, url, alt_text) VALUES ($1, $2, $3)",
		image.Title, image.URL, image.AltText,
	)
	if err != nil {
		return fmt.Errorf("unable to insert an entry: %v", err)
	}
	return nil
}

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

	http.HandleFunc("/images.json", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		indent, err := strconv.Atoi(queryParams.Get("indent"))
		if err != nil {
			http.Error(w, "Indent should be an integer", http.StatusBadRequest)
		}

		if r.Method == "GET" {
			images, err := fetchImages(conn)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to fetch images: %v", err)
				os.Exit(1)
			}

			b, err := json.MarshalIndent(images, "", strings.Repeat(" ", indent))
			if err != nil {
				http.Error(w, "Cannot serialize object", http.StatusInternalServerError)
			}

			w.Header().Add("Content-Type", "text/json")
			w.Write([]byte(b))
		} else if r.Method == "POST" {
			var image Image
			body, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot read request body: %v", err)
				http.Error(w, "Cannot read request body", http.StatusInternalServerError)
				return
			}
			if err := json.Unmarshal(body, &image); err != nil {
				fmt.Fprintf(os.Stderr, "cannot deserialize image: %v", err)
				http.Error(w, "Cannot deserialize image", http.StatusInternalServerError)
				return
			}
			if err := addImage(conn, image); err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}

			b, err := json.MarshalIndent(image, "", strings.Repeat(" ", indent))
			if err != nil {
				http.Error(w, "Cannot serialize object", http.StatusInternalServerError)
			}
			w.Header().Add("Content-Type", "text/json")
			w.Write([]byte(b))
		}
	})

	http.ListenAndServe(":8080", nil)
}
