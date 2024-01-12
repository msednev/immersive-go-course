package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
)

type DbConfig struct {
	DbUrl string
	Port  int
}

func Run(config DbConfig) error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("unable to create to database connection: %w", err)
	}
	defer conn.Close(context.Background())

	mux := http.NewServeMux()

	mux.HandleFunc("/images.json", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.EscapedPath())
		indent := r.URL.Query().Get("indent")

		var response []byte
		var responseErr error

		if r.Method == "GET" {
			images, err := FetchImages(conn)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
			response, responseErr = MarschalWithIndent(images, indent)

		} else if r.Method == "POST" {
			image, err := AddImage(conn, r)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
			response, responseErr = MarschalWithIndent(image, indent)
		}

		if responseErr != nil {
			fmt.Fprintln(os.Stderr, responseErr.Error())
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8082")
		if _, err := w.Write(response); err != nil {
			fmt.Fprintf(os.Stderr, "unable to write response: %v", err)
		}
	})

	log.Printf("port: %v", config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Port), mux)

}
