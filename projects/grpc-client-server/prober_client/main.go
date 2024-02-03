// Package main implements a client for Prober service.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/CodeYourFuture/immersive-go-course/grpc-client-server/prober"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	var endpoint string
	var n int64

	flag.StringVar(&endpoint, "endpoint", "http://www.google.com", "endpoint to make request to")
	flag.Int64Var(&n, "n", 1, "number of requests")
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProberClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second) // TODO: add a timeout
	defer cancel()

	// TODO: endpoint should be a flag
	// TODO: add number of times to probe
	r, err := c.DoProbes(ctx, &pb.ProbeRequest{Endpoint: endpoint, NRequests: n})
	if err != nil {
		log.Fatalf("could not probe: %v", err)
	}
	log.Printf("Average Response Time: %f", r.GetMeanLatencyMsecs())
}
