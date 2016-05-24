package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/frrakn/treebeer/util/config"
	"github.com/frrakn/treebeer/util/handle"
	"github.com/juju/errors"
)

type configuration struct {
	Season int32
	Period string
}

type riotSeason struct {
	SeasonName  string
	SeasonSplit string
	ProTeams    []riotTeam
	ProPlayers  []riotPlayer
}

type riotTeam struct {
	Id        int32
	RiotId    int32
	Name      string
	ShortName string
}

type riotPlayer struct {
	Id        int32
	RiotId    int32
	Name      string
	ProTeamId int32
}

var (
	season      int32
	conf        configuration
	checkPeriod time.Duration
	address     = "http://fantasy.na.lolesports.com/en-US/api/season/%d"
)

func main() {
	initialize()

	var lastRequest []byte

	for {
		time.Sleep(checkPeriod)

		s := getSeason()

		var seasonCtx riotSeason
		b := requestSeasonData(s)
		if bytes.Compare(b, lastRequest) != 0 {
			seasonCtx = parseSeasonData(b)
			lastRequest = b
		}

		// TODO(fchen): convert to messages and make RPC call to ContextManager
		fmt.Println(seasonCtx)
	}
}

func initialize() {
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

func getSeason() int32 {
	return conf.Season
}

func requestSeasonData(s int32) []byte {
	res, err := http.Get(fmt.Sprintf(address, s))
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

func parseSeasonData(body []byte) riotSeason {
	var season riotSeason
	err := json.Unmarshal(body, &season)
	if err != nil {
		handle.Error(err)
	}

	return season
}
