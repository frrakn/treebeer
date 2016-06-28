package db

import (
	"time"

	ctxPb "github.com/frrakn/treebeer/context/proto"
)

func (g *Game) ToPB() *ctxPb.SavedGame {
	return &ctxPb.SavedGame{
		Game: &ctxPb.Game{
			Lcsid:       int32(g.LcsID),
			Riotgameid:  g.RiotGameID,
			Riotmatchid: g.RiotMatchID,
			Redteamid:   int32(g.RedTeamID),
			Blueteamid:  int32(g.BlueTeamID),
			Gamestart:   g.GameStart.UTC().Unix(),
			Gameend:     g.GameEnd.UTC().Unix(),
		},
		Gameid: int32(g.GameID),
	}
}

func (g *Game) FromPB(game *ctxPb.Game, id GameID) {
	g.GameID = id
	g.LcsID = LcsID(game.Lcsid)
	g.RiotGameID = game.Riotgameid
	g.RiotMatchID = game.Riotmatchid
	g.BlueTeamID = TeamID(game.Blueteamid)
	g.RedTeamID = TeamID(game.Redteamid)
	g.GameStart = time.Unix(game.Gamestart, 0)
	g.GameEnd = time.Unix(game.Gameend, 0)
}

func (g *Game) FromSavedPB(game *ctxPb.SavedGame) {
	g.GameID = GameID(game.Gameid)
	g.LcsID = LcsID(game.Game.Lcsid)
	g.RiotGameID = game.Game.Riotgameid
	g.RiotMatchID = game.Game.Riotmatchid
	g.RedTeamID = TeamID(game.Game.Redteamid)
	g.BlueTeamID = TeamID(game.Game.Blueteamid)
	g.GameStart = time.Unix(game.Game.Gamestart, 0)
	g.GameEnd = time.Unix(game.Game.Gameend, 0)
}

func (g *Game) EqualsPB(game *ctxPb.Game) bool {
	return (game.Lcsid == int32(g.LcsID) &&
		game.Riotgameid == g.RiotGameID &&
		game.Riotmatchid == g.RiotMatchID &&
		game.Redteamid == int32(g.RedTeamID) &&
		game.Blueteamid == int32(g.BlueTeamID) &&
		game.Gamestart == g.GameStart.Unix() &&
		game.Gameend == g.GameEnd.Unix())
}

func (g *Game) SetIDs(red TeamID, blue TeamID) {
	g.RedTeamID = red
	g.BlueTeamID = blue
}
