package main

import (
	"flag"
	"fmt"
	"os"
	"servers/static"
)

func main() {
	pathPtr := flag.String("path", "", "where to read static files from")
	portPtr := flag.Int("port", -1, "a port this server will be running on")
	flag.Parse()

	if *pathPtr == "" {
		fmt.Fprintf(os.Stderr, "path is not set, exiting...")
		os.Exit(1)
	}

	if *portPtr == -1 {
		fmt.Fprintf(os.Stderr, "server port is not set, exiting...")
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, static.Run(static.Config{Dir: *pathPtr, Port: *portPtr}))
}
