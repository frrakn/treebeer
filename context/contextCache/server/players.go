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

func (p *players) all() []*db.Player {
	ps := []*db.Player{}
	p.RLock()
	for _, player := range p.m {
		ps = append(ps, player)
	}
	p.RUnlock()
	return ps
}
