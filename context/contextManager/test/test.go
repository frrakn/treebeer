package main

import (
	"fmt"
	"log"

	pb "github.com/frrakn/treebeer/context/proto"
	"github.com/juju/errors"
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
	c := pb.NewLiveStatUpdateClient(conn)

	/*updates := &pb.BatchUpdates{
		&pb.Teams{
			[]*pb.Team{
				&pb.Team{0, 1, 1, "harro", "HRO"},
				&pb.Team{0, 2, 2, "test", "TST"},
			},
		},
		&pb.Teams{},
		&pb.Players{
			[]*pb.Player{
			//				&pb.Player{0, 1, 1, "herp", 1},
			//				&pb.Player{0, 2, 2, "derp", 1},
			//				&pb.Player{0, 3, 3, "doo", 2},
			},
		},
		&pb.Players{},
		&pb.Games{
			[]*pb.Game{
			//				&pb.Game{0, 1, "a", "a", 1, 2, time.Now().Unix(), 0},
			},
		},
		&pb.Games{},
		&pb.Stats{
			[]*pb.Stat{
			//				&pb.Stat{0, "kills"},
			//				&pb.Stat{1, "deaths"},
			},
		},
		&pb.Stats{},
	}*/

	// Contact the server and print out its response.
	team, err := c.GetTeam(context.Background(), &pb.Team{1, 1, "a", "a"})

	fmt.Println(team)
	fmt.Println(errors.ErrorStack(err))

	stat, err := c.GetStat(context.Background(), &pb.Stat{"kills"})

	fmt.Println(stat)
	fmt.Println(errors.ErrorStack(err))
}
