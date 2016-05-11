package main

import (
	"fmt"
	"log"
	"time"

	pb "github.com/frrakn/treebeer/contextManager/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("104.196.125.93:9321", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewContextUpdateClient(conn)

	updates := &pb.BatchUpdates{
		&pb.Teams{
			[]*pb.Team{
				&pb.Team{0, 1, 1, "harro", "HRO"},
				&pb.Team{0, 2, 2, "test", "TST"},
			},
		},
		&pb.Teams{},
		&pb.Players{
			[]*pb.Player{
				&pb.Player{0, 1, 1, "herp", 1},
				&pb.Player{0, 2, 2, "derp", 1},
				&pb.Player{0, 3, 3, "doo", 2},
			},
		},
		&pb.Players{},
		&pb.Games{
			[]*pb.Game{
				&pb.Game{0, 1, "a", "a", 1, 2, time.Now().Unix(), 0},
			},
		},
		&pb.Games{},
	}

	// Contact the server and print out its response.
	_, err = c.BatchUpdate(context.Background(), updates)

	fmt.Println(err)
}
