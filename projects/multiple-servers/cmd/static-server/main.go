package main

import (
	"flag"
	"servers/static"
)

func main() {
	pathPtr := flag.String("path", "", "where to read static files from")
	portPtr := flag.Int("port", 0, "a port this server will be running on")
	flag.Parse()
	static.Run(*pathPtr, *portPtr)
}