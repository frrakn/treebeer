package main

import (
	"flag"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/frrakn/treebeer/context/contextManager/manager"
	"github.com/frrakn/treebeer/context/db"
	ctxPb "github.com/frrakn/treebeer/context/proto"
	"github.com/frrakn/treebeer/util/config"
	"github.com/frrakn/treebeer/util/handle"

	"github.com/juju/errors"
)

type configuration struct {
	DB       string
	Port     string
	Keyfiles db.Keyfiles
}

const (
	tlsProfile = "contextManager"
)

var (
	conf      configuration
	ctxServer *manager.Server
)

func main() {
	serveRpc(conf.Port)
}

func init() {
	flag.Parse()

	ctxServer = manager.NewServer()

	err := config.LoadConfig(&conf)
	if err != nil {
		handle.Fatal(errors.Annotate(err, "Failed to load configuration"))
	}

	ctxServer.SqlDB, err = db.InitDB(conf.DB+tlsProfile, tlsProfile, conf.Keyfiles)
	if err != nil {
		handle.Fatal(errors.Annotate(err, "Failed to load DB"))
	}

	season, err := db.GetSeasonContext(ctxServer.SqlDB)
	if err != nil {
		handle.Fatal(errors.Annotate(err, "Failed to load season data from DB"))
	}

	ctxServer.Initialize(season)
}

func serveRpc(port string) {
	l, err := net.Listen("tcp", port)
	if err != nil {
		handle.Fatal(errors.Annotate(err, fmt.Sprintf("Unable to get listener on port %d", port)))
	}

	rpcserv := grpc.NewServer()
	ctxPb.RegisterSeasonUpdateServer(rpcserv, ctxServer)
	rpcserv.Serve(l)
}
