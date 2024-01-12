package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jackc/pgx/v4"
)

type Image struct {
	Title   string `json:"title"`
	AltText string `json:"alt_text"`
	Url     string `json:"url"`
}

func FetchImages(conn *pgx.Conn) ([]Image, error) {
	var images []Image
	var title, url, altText string

	rows, err := conn.Query(context.Background(), "SELECT title, url, alt_text FROM public.images")
	if err != nil {
		return nil, fmt.Errorf("unable to execute query: %w", err)
	}

	for rows.Next() {
		if err = rows.Scan(&title, &url, &altText); err != nil {
				return nil, fmt.Errorf("unable to extract columns: %w", err)
		}
		images = append(images, Image{Title: title, Url: url, AltText: altText})
	}
	rows.Close()

	return images, nil
}

func AddImage(conn *pgx.Conn, r *http.Request) (Image, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return Image{}, fmt.Errorf("unable to read request body: %w", err)
	}

	var image Image
	err = json.Unmarshal(body, &image)
	if err != nil {
		return Image{}, fmt.Errorf("unable to deserialize body: %w", err)
	}

	_, err = conn.Exec(
		context.Background(),
		"INSERT INTO public.images (title, url, alt_text) VALUES ($1, $2, $3)",
		image.Title, image.Url, image.AltText,
	)
	if err != nil {
		return Image{}, fmt.Errorf("unable to insert an entry: %w", err)
	}
	return image, nil
}