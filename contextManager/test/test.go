package main

import (
	"log"

	pb "github.com/frrakn/treebeer/contextManager/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewContextUpdateClient(conn)

	// Contact the server and print out its response.
	_, _ = c.CreateTeam(context.Background(), &pb.Team{0, 1, 1, "harro", "HRO"})
}
