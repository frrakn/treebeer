package main

import (
	"flag"
	"os"
	"os/signal"

	"golang.org/x/sys/unix"

	"github.com/frrakn/treebeer/liveData/translator/ws"
	"github.com/frrakn/treebeer/util/config"
	"github.com/frrakn/treebeer/util/handle"
)

type configuration struct {
	Address       string
	Path          string
	SocketOptions map[string]string
}

var (
	conf configuration
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

	listener := ws.NewListener(conf.Address, conf.Path, conf.SocketOptions)
	listener.Start()
	defer listener.Stop()

	<-sigs
}
