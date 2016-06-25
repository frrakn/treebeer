package contextStore

import (
	"sync"

	"github.com/frrakn/treebeer/context/db"
)

type players struct {
	m map[db.LcsID]*sync.RWMutex
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
