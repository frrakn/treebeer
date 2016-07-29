package main

import (
	"flag"
	"os"
	"os/signal"

	"golang.org/x/sys/unix"

	"github.com/frrakn/treebeer/liveData/translator/contextStore"
	"github.com/frrakn/treebeer/liveData/translator/server"
	"github.com/frrakn/treebeer/liveData/translator/ws"
	"github.com/frrakn/treebeer/util/config"
	"github.com/frrakn/treebeer/util/handle"
)

type configuration struct {
	Listener     *ws.Configuration
	ContextStore *contextStore.Configuration
	Server       *server.Configuration
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

	ctxStore := contextStore.New(conf.ContextStore)
	listener := ws.NewListener(conf.Listener)

	server := server.NewServer(conf.Server, ctxStore, listener)

	server.Start()

	defer server.Stop()

	for {
		select {
		case err := <-server.Errors:
			handle.Error(err)
		case <-sigs:
			return
		}
	}
}
