package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/frrakn/treebeer/context/contextUpdater/schema"
	"github.com/frrakn/treebeer/util/config"
	"github.com/frrakn/treebeer/util/handle"
	"github.com/juju/errors"

	ctxPb "github.com/frrakn/treebeer/context/proto"
)

type configuration struct {
	Season                int32
	Period                string
	ContextManagerAddress string
}

var (
	season      int32
	conf        configuration
	checkPeriod time.Duration
	lcsAddress  = "http://fantasy.na.lolesports.com/en-US/api/season/%d"
)

func main() {
	conn, err := grpc.Dial(conf.ContextManagerAddress, grpc.WithInsecure())
	if err != nil {
		handle.Fatal(errors.Annotatef(err, "Failed to connect to rpc server"))
	}
	defer closeConn(conn)

	var lastRequest []byte
	client := ctxPb.NewSeasonUpdateClient(conn)

	for {
		time.Sleep(checkPeriod)

		s := getSeason()

		var seasonCtx schema.RiotSeason
		b := requestSeasonData(s)
		if bytes.Compare(b, lastRequest) != 0 {
			seasonCtx = parseSeasonData(b)
			lastRequest = b

			seasonUpdate, err := seasonCtx.ToSeasonUpdates()
			if err != nil {
				handle.Error(errors.Trace(err))
				continue
			}
			client.SeasonUpdate(context.Background(), seasonUpdate)
		}
	}
}

func init() {
	flag.Parse()

	err := config.LoadConfig(&conf)
	if err != nil {
		handle.Fatal(err)
	}

	checkPeriod, err = time.ParseDuration(conf.Period)
	if err != nil {
		handle.Fatal(errors.Annotatef(err, "Unable to parse ContextUpdater check period of %s", conf.Period))
	}
}

func closeConn(conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {
		handle.Fatal(errors.Trace(err))
	}
}

func getSeason() int32 {
	return conf.Season
}

func requestSeasonData(s int32) []byte {
	res, err := http.Get(fmt.Sprintf(lcsAddress, s))
	if err != nil {
		handle.Error(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		handle.Error(err)
	}

	return body
}

func parseSeasonData(body []byte) schema.RiotSeason {
	var season schema.RiotSeason
	err := json.Unmarshal(body, &season)
	if err != nil {
		handle.Error(err)
	}

	return season
}
