package main

import (
	"flag"
	"time"

	"github.com/frrakn/treebeer/context/contextServer/server"
	"github.com/frrakn/treebeer/context/db"
	"github.com/frrakn/treebeer/util/config"
	"github.com/frrakn/treebeer/util/handle"
	"github.com/juju/errors"
)

type configuration struct {
	DB       string
	Port     string
	Keyfiles db.Keyfiles
	Period   string
}

const (
	tlsProfile = "contextServer"
)

var (
	conf        configuration
	checkPeriod time.Duration
	ctxServer   *server.Server
)

func main() {
	forever := make(chan struct{})
	ctxServer.Run(conf.Port)
	defer ctxServer.Stop()
	<-forever
}

func init() {
	flag.Parse()

	err := config.LoadConfig(&conf)
	if err != nil {
		handle.Fatal(errors.Annotate(err, "Failed to load configuration"))
	}

	sqlDB, err := db.InitDB(conf.DB+tlsProfile, tlsProfile, conf.Keyfiles)
	if err != nil {
		handle.Fatal(errors.Annotate(err, "Failed to load DB"))
	}

	checkPeriod, err = time.ParseDuration(conf.Period)
	if err != nil {
		handle.Fatal(errors.Annotatef(err, "Unable to parse ContextUpdater check period of %s", conf.Period))
	}

	ctxServer = server.NewServer(sqlDB)
}
