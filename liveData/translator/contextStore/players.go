package contextStore

import (
	"sync"

	"github.com/frrakn/treebeer/context/db"
)

type players struct {
	m map[db.RiotID]*db.Player
	sync.RWMutex
}

func newPlayers() *players {
	return &players{
		m: make(map[db.RiotID]*db.Player),
	}
}

func (p *players) batchUpdate(players []*db.Player) {
	p.Lock()
	for _, player := range players {
		p.m[player.RiotID] = player
	}
	p.Unlock()
}

func (p *players) get(id db.RiotID) (player *db.Player, ok bool) {
	p.RLock()
	player, ok = p.m[id]
	p.RUnlock()
	return
}

func (p *players) set(id db.RiotID, player *db.Player) {
	p.RLock()
	p.m[id] = player
	p.RUnlock()
}
