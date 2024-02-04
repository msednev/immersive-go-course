package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	pb "github.com/CodeYourFuture/immersive-go-course/grpc-client-server/prober"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement prober.ProberServer.
type server struct {
	pb.UnimplementedProberServer
}

func (s *server) DoProbes(ctx context.Context, in *pb.ProbeRequest) (*pb.ProbeReply, error) {
	// TODO: support a number of repetitions and return average latency
	var sumElapsedMsecs float32 = 0.0
	for i := 0; i < int(in.GetNRequests()); i++ {
		start := time.Now()
		resp, err := http.Get(in.GetEndpoint())	// TODO: add error handling here and check the response code
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("non-ok response status")
		}
		elapsed := time.Since(start)
		sumElapsedMsecs += float32(elapsed / time.Millisecond)
	}

	return &pb.ProbeReply{MeanLatencyMsecs: sumElapsedMsecs / float32(in.GetNRequests())}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProberServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
