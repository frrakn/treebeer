package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"golang.org/x/sys/unix"

	"github.com/frrakn/treebeer/liveData/translator/contextStore"
	"github.com/frrakn/treebeer/liveData/translator/ws"
	"github.com/frrakn/treebeer/liveData/translator/ws/schema"
	"github.com/frrakn/treebeer/util/config"
	"github.com/frrakn/treebeer/util/handle"
)

type configuration struct {
	Listener     *ws.Configuration
	ContextStore *contextStore.Configuration
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

	go handleErrors(ctxStore, listener)

	ctxStore.Start()
	listener.Start()
	defer ctxStore.Stop()
	defer listener.Stop()

	for {
		select {
		case liveStats := <-listener.Stats:
			handleStats(ctxStore, liveStats)
		case err := <-listener.Errors:
			handle.Error(err)
		case <-sigs:
			return
		}
	}
}

func handleStats(ctxStore *contextStore.ContextStore, stats *schema.LiveStats) {
	fmt.Println(stats)
}

func handleErrors(ctxStore *contextStore.ContextStore, listener *ws.Listener) {
	for {
		select {
		case err := <-ctxStore.Errors:
			handle.Error(err)
		case err := <-listener.Errors:
			handle.Error(err)
		}
	}
}
