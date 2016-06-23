package manager

import (
	"sync"

	"github.com/frrakn/treebeer/context/db"
)

type games struct {
	m map[db.LcsID]*db.Game
	sync.RWMutex
}

func (g *games) batchUpdate(games []*db.Game) {
	g.Lock()
	for _, game := range games {
		g.m[game.LcsID] = game
	}
	g.Unlock()
}

func (g *games) get(id db.LcsID) (game *db.Game, ok bool) {
	g.RLock()
	game, ok = g.m[id]
	g.RUnlock()
	return
}
