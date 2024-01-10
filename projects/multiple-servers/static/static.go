package static

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

type Config struct {
	Dir string
	Port int
}

func Run(config Config) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(config.Dir, r.URL.EscapedPath())
		log.Println(r.Method, r.URL.EscapedPath(), path)
		http.ServeFile(w, r, path)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", config.Port), mux)
}