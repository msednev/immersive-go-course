package main

import (
	"fmt"
	"go-ls/cmd"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "wrong number of arguments\n")
		os.Exit(1)
	}

	dirname := os.Args[1]
	
	if err:= cmd.Execute(dirname); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	os.Exit(0)
	
}
