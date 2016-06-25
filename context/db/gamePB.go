package db

import (
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
			Gamestart:   g.GameStart,
			Gameend:     g.GameEnd,
		},
		Gameid: int32(g.GameID),
	}
}

func (g *Game) FromPB(game *ctxPb.Game, id GameID) {
	g.GameID = id
	g.LcsID = LcsID(game.Lcsid)
	g.RiotGameID = game.Riotgameid
	g.RiotMatchID = game.Riotmatchid
	// RedTeamID foreign key omitted for internal translation
	// BlueTeamID foreign key omitted for internal translation
	g.GameStart = game.Gamestart
	g.GameEnd = game.Gameend
}

func (g *Game) FromSavedPB(game *ctxPb.SavedGame) {
	g.GameID = GameID(game.Gameid)
	g.LcsID = LcsID(game.Game.Lcsid)
	g.RiotGameID = game.Game.Riotgameid
	g.RiotMatchID = game.Game.Riotmatchid
	g.RedTeamID = TeamID(game.Game.Redteamid)
	g.BlueTeamID = TeamID(game.Game.Blueteamid)
	g.GameStart = game.Game.Gamestart
	g.GameEnd = game.Game.Gameend
}

func (g *Game) SetIDs(red TeamID, blue TeamID) {
	g.RedTeamID = red
	g.BlueTeamID = blue
}
