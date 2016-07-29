/*
====================
LIVE DATA AGGREGATOR
====================
Provides data aggregation service for incoming map[string]int

*/

package main

import (
	"flag"
	"os"
	"os/signal"

	"golang.org/x/sys/unix"

	"github.com/frrakn/treebeer/liveData/aggregator/server"
	"github.com/frrakn/treebeer/util/config"
	"github.com/frrakn/treebeer/util/handle"
)

type configuration struct {
	Server *server.Configuration
}

var (
	conf                   configuration
	promotionBroadcastSema = make(chan struct{}, 1)
)

func init() {
	flag.Parse()

	err := config.LoadConfig(&conf)
	if err != nil {
		handle.Fatal(err)
	}
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, unix.SIGINT, unix.SIGTERM)

	s := server.NewServer(conf.Server)
	s.Start()

	for {
		select {
		case err := <-s.Errors:
			handle.Error(err)
		case <-sigs:
			return
		}
	}
}
