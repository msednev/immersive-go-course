package main

import (
	"flag"
	"fmt"
	"os"
	"servers/api"
)

func main() {
	dbUrl := os.Getenv("DATABASE_URL")
	portPtr := flag.Int("port", -1, "port the server listens on")
	flag.Parse()

	if dbUrl == "" {
		fmt.Fprintf(os.Stderr, "DATABASE_URL is not set, exiting...")
		os.Exit(1)
	}

	if *portPtr == -1 {
		fmt.Fprintf(os.Stderr, "server port is not set, exiting...")
		os.Exit(1)
	}

	api.Run(api.DbConfig{DbUrl: dbUrl, Port: *portPtr})
}
