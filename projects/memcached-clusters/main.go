package main

import (
	"flag"
	"log"
	"strconv"
	"strings"
)


func parseInts(s string) ([]int, error) {
	var ints []int
	for _, ss := range strings.Split(s, ",") {
		i, err := strconv.Atoi(ss)
		if err != nil {
			return nil, err
		}
		ints = append(ints, i)
	}
	return ints, nil
}

func main() {
	var mcrouter int
	var memcachedStr string

	flag.IntVar(&mcrouter, "mcrouter", 0, "mcrouter port")
	flag.StringVar(&memcachedStr,"memcacheds", "", "comma-separated list of memcached ports")
	memcacheds, err := parseInts(memcachedStr)
	if err != nil {
		log.Fatalf("failed to parse command line arguments: %v: %v", memcachedStr, err)
	}

	
}
