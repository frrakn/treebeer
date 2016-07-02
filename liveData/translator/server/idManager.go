package server

import (
	"sync"

	"github.com/frrakn/treebeer/context/db"
)

type game struct {
	sync.RWMutex
	gameID  db.GameID
	teams   map[string]db.TeamID
	players map[string]db.PlayerID
}

func NewGame() *game {
	return &game{
		teams:   make(map[string]db.TeamID),
		players: make(map[string]db.PlayerID),
	}
}

func (g *game) getID() db.GameID {
	g.RLock()
	id := g.gameID
	g.RUnlock()
	return id
}

func (g *game) getTeam(idStr string) (db.TeamID, bool) {
	g.RLock()
	id, ok := g.teams[idStr]
	g.RUnlock()
	return id, ok
}

func (g *game) getPlayer(idStr string) (db.PlayerID, bool) {
	g.RLock()
	id, ok := g.players[idStr]
	g.RUnlock()
	return id, ok
}

func (g *game) setID(id db.GameID) {
	g.Lock()
	g.gameID = id
	g.Unlock()
}

func (g *game) setTeam(idStr string, id db.TeamID) {
	g.Lock()
	g.teams[idStr] = id
	g.Unlock()
}

func (g *game) setPlayer(idStr string, id db.PlayerID) {
	g.Lock()
	g.players[idStr] = id
	g.Unlock()
}
