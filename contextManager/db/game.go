package db

import "time"

type Game struct {
	GameId      int32
	LcsId       int32
	RiotGameId  string
	RiotMatchId string
	RedTeamId   int32
	BlueTeamId  int32
	GameStart   time.Time
	GameEnd     time.Time
}
