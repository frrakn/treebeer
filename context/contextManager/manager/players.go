package manager

import (
	"sync"

	"github.com/frrakn/treebeer/context/db"
)

type players struct {
	m map[db.LcsID]*db.Player
	sync.RWMutex
}

func (p *players) batchUpdate(players []*db.Player) {
	p.Lock()
	for _, player := range players {
		p.m[player.LcsID] = player
	}
	p.Unlock()
}

func (p *players) batchUpdateWithID(players []*playerAndTeamID) {
	p.Lock()
	for _, pWithID := range players {
		player := pWithID.Player
		p.m[player.LcsID] = player
	}
	p.Unlock()
}

func (p *players) get(id db.LcsID) (player *db.Player, ok bool) {
	p.RLock()
	player, ok = p.m[id]
	p.RUnlock()
	return
}

func (p *players) set(id db.LcsID, player *db.Player) {
	p.RLock()
	p.m[id] = player
	p.RUnlock()
}
