package contextStore

import (
	"sync"

	"github.com/frrakn/treebeer/context/db"
)

type games struct {
	m map[db.LcsID]*db.Game
	sync.RWMutex
}

func newGames() *games {
	return &games{
		m: make(map[db.LcsID]*db.Game),
	}
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

func (g *games) set(id db.LcsID, game *db.Game) {
	g.RLock()
	g.m[id] = game
	g.RUnlock()
}
