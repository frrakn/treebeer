package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/frrakn/treebeer/util/config"
	"github.com/frrakn/treebeer/util/handle"
)

type configuration struct {
	Season int32
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
	conf    configuration
	address = "http://fantasy.na.lolesports.com/en-US/api/season/%d"
)

func main() {
	flag.Parse()

	err := config.LoadConfig(&conf)
	if err != nil {
		handle.Fatal(err)
	}

	s := getSeason()
	b := requestSeason(s)

	fmt.Println(b)
}

func getSeason() int32 {
	return conf.Season
}

func requestSeason(s int32) (seasonCtx riotSeason) {
	res, err := http.Get(fmt.Sprintf(address, s))
	if err != nil {
		handle.Error(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		handle.Error(err)
	}

	var season riotSeason
	err = json.Unmarshal(body, &season)
	if err != nil {
		handle.Error(err)
	}

	return season
}
