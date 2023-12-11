package main

import (
	"net/http"
	"fmt"
	"os"
	"strconv"
	"time"
	"io"
)

const address = "http://localhost:8080"

func main() {
	resp, err := http.Get(address)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 429 {
		timeToSleep, err := strconv.Atoi(resp.Header.Get("Retry-After"))
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"failed to determine how long to sleep: %v\n, waiting 5 s\n",
				err,
			)
			time.Sleep(5 * time.Second)
			main()
		}

		if timeToSleep > 1 {
			fmt.Fprintf(
				os.Stderr,
				"Things are running a bit slow, waiting %v seconds\n",
				timeToSleep,
			)
			time.Sleep(time.Duration(timeToSleep) * time.Second)
			main()
		}

		if timeToSleep > 5 {
			fmt.Fprintln(
				os.Stderr,
				"Cannot get weather",
			)
			os.Exit(2)
		}
	}
	
	if resp.StatusCode == 200 {
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't read response body: %v\n", err)
		}
		fmt.Fprintln(os.Stdout, string(msg))
		os.Exit(0)
	}

	fmt.Fprintln(os.Stderr, "Can't connect to server")
	os.Exit(3)
}
