package main

import (
	"flag"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
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

func getMemcachedClients(server string, ports []int) []*memcache.Client {
	mcClients := make([]*memcache.Client, len(ports))
	for index, port := range ports {
		mcClients[index] = memcache.New(fmt.Sprintf("%s:%d", server, port))
	}
	return mcClients
}

func main() {
	var mcrouterPort int
	var memcachedPorts string

	flag.IntVar(&mcrouterPort, "mcrouter", 0, "mcrouter port")
	flag.StringVar(&memcachedPorts, "memcacheds", "", "comma-separated list of memcached ports")
	flag.Parse()
	memcachedPortsParsed, err := parseInts(memcachedPorts)
	if err != nil {
		log.Fatalf("failed to parse command line arguments: %v: %v", memcachedPorts, err)
	}

	mcrouterClient := memcache.New(fmt.Sprintf("127.0.0.1:%d", mcrouterPort))
	memcachedClients := getMemcachedClients("127.0.0.1", memcachedPortsParsed)

	const key = "mykey"
	const value = "my data"

	if err := mcrouterClient.Set(&memcache.Item{Key: key, Value: []byte(value)}); err == nil {
		log.Fatalf("failed to set key: %v", err)
	}

	var wg sync.WaitGroup

	values := make(chan []byte, len(memcachedClients))

	for _, client := range memcachedClients {
		wg.Add(1)
		go func(client *memcache.Client) {
			defer wg.Done()
			item, err := client.Get(key)
			if err != nil {
				log.Printf("cannot retrieve item: %v", err)
			}
			values <- item.Value
		}(client)
	}
	wg.Wait()
	close(values)

	for el := range values {
		if !reflect.DeepEqual([]byte(value), el) {
			fmt.Println("the caches operate in sharded mode")
			return
		}
	}

	fmt.Println("the caches operate in replicated mode")
}
