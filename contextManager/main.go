package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"net"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/frrakn/treebeer/contextManager/db"
	ctxPb "github.com/frrakn/treebeer/contextManager/proto"
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

type server struct {
	sqldb *sqlx.DB
}

func main() {
	flag.Parse()

	var c configuration
	err := config.LoadConfig(&c)
	if err != nil {
		handle.Fatal(errors.Annotate(err, "Failed to load configuration"))
	}
	liveStatDB := initDB(c.DB, c.Keyfiles)
	season, err := db.GetSeasonContext(liveStatDB)

	serveRpc(c.Port, liveStatDB)
}

func errorString(err error) string {
	return fmt.Sprintf("%s", err)
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

	liveStatDB, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		handle.Fatal(errors.Annotate(err, fmt.Sprintf("Unable to connect to database at %s", dsn)))
	}

	return liveStatDB
}

func serveRpc(port string, liveStatDB *sqlx.DB) {
	l, err := net.Listen("tcp", port)
	if err != nil {
		handle.Fatal(errors.Annotate(err, fmt.Sprintf("Unable to get listener on port %d", port)))
	}

	rpcserv := grpc.NewServer()
	ctxPb.RegisterSeasonUpdateServer(rpcserv, &server{liveStatDB})
	rpcserv.Serve(l)
}

func (s *server) SeasonUpdate(ctx context.Context, updates *ctxPb.SeasonUpdates) (*ctxPb.Empty, error) {
	return &ctxPb.Empty{}, nil
}
