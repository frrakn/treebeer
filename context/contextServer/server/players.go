package server

import (
	"sync"

	"github.com/frrakn/treebeer/context/db"
)

type players struct {
	m map[db.PlayerID]*db.Player
	sync.RWMutex
}

func (p *players) batchUpdate(players []*db.Player) {
	p.Lock()
	for _, player := range players {
		p.m[player.PlayerID] = player
	}
	p.Unlock()
}

func (p *players) get(id db.PlayerID) (player *db.Player, ok bool) {
	p.RLock()
	player, ok = p.m[id]
	p.RUnlock()
	return
}

func (p *players) getForTeam(id db.TeamID) []*db.Player {
	players := make([]*db.Player, 0)

	p.RLock()
	for _, player := range p.m {
		if player.TeamID == id {
			players = append(players, player)
		}
	}
	p.RUnlock()
	return players
}
