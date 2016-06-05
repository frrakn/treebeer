package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"net"

	"google.golang.org/grpc"

	"github.com/frrakn/treebeer/context/contextManager/manager"
	"github.com/frrakn/treebeer/context/db"
	ctxPb "github.com/frrakn/treebeer/context/proto"
	"github.com/frrakn/treebeer/util/config"
	"github.com/frrakn/treebeer/util/handle"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
)

type configuration struct {
	DB       string
	Port     string
	Keyfiles keyfiles
}

type keyfiles struct {
	CaCert     string
	ClientCert string
	ClientKey  string
}

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
	ctxServer.SqlDB = initDB(conf.DB, conf.Keyfiles)

	season, err := db.GetSeasonContext(ctxServer.SqlDB)
	if err != nil {
		handle.Fatal(errors.Annotate(err, "Failed to load season data from DB"))
	}

	ctxServer.Initialize(season)
}

func initDB(dsn string, keys keyfiles) *sqlx.DB {
	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile(keys.CaCert)
	if err != nil {
		handle.Fatal(errors.Annotate(err, fmt.Sprintf("Unable to access database credentials at %s", keys.CaCert)))
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		handle.Fatal(errors.Annotate(err, "Unabe to append PEM."))
	}

	clientCert := make([]tls.Certificate, 0, 1)
	certs, err := tls.LoadX509KeyPair(keys.ClientCert, keys.ClientKey)
	if err != nil {
		handle.Fatal(errors.Annotate(err, fmt.Sprintf("Unable to access database credentials at %s and %s", keys.ClientCert, keys.ClientKey)))
	}
	clientCert = append(clientCert, certs)

	mysql.RegisterTLSConfig("treebeer", &tls.Config{
		RootCAs:            rootCertPool,
		Certificates:       clientCert,
		InsecureSkipVerify: true,
	})

	sqldb, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		handle.Fatal(errors.Annotate(err, fmt.Sprintf("Unable to connect to database at %s", dsn)))
	}

	return sqldb
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
