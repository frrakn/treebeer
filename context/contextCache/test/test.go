package main

import (
	"fmt"
	"log"

	pb "github.com/frrakn/treebeer/context/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server.
	//conn, err := grpc.Dial("104.196.125.93:9321", grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSeasonContextClient(conn)

	// Contact the server and print out its response.
	fullCtx, err := c.GetContext(context.Background(), &pb.Empty{})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(fullCtx)
}
