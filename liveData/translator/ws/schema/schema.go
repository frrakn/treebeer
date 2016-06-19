package schema

import "github.com/frrakn/treebeer/context/db"

type LiveStats struct {
	Games map[db.LcsID]Game
}

type Game struct {
	TeamStats   map[int32]Team
	PlayerStats map[int32]Player
}

type Team struct {
	TeamID           int32
	BaronsKilled     int32
	DragonsKilled    int32
	FirstBlood       bool
	TowersKilled     int32
	InhibitorsKilled int32
	MatchVictory     int32
	MatchDefeat      int32
	Color            string
}

type Player struct {
	ParticipantID  int32
	TeamID         int32
	SummonerName   string
	ChampionName   string
	ChampionID     int32
	SkinIndex      int32
	ProfileIconID  int32
	SummonerSpell1 int32
	SummonerSpell2 int32
	Kills          int32
	Deaths         int32
}
