package main

import (
	"fmt"

	"github.com/frrakn/treebeer/context/db"
	"github.com/frrakn/treebeer/context/poller"
)

func main() {
	p, _ := poller.NewRPCPoller("localhost:8080")

	var (
		season *db.SeasonContext
		err    error
	)

	go p.Run()
	defer p.Stop()

	for {
		select {
		case season = <-p.Updates:
			fmt.Println(season.Players[0])
		case err = <-p.Errors:
			fmt.Println(err)
		}
	}
}
