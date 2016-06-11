package server

import (
	"sync"

	"github.com/frrakn/treebeer/context/db"
)

type teams struct {
	m map[db.TeamID]*db.Team
	sync.RWMutex
}

func (t *teams) batchUpdate(teams []*db.Team) {
	t.Lock()
	for _, team := range teams {
		t.m[team.TeamID] = team
	}
	t.Unlock()
}

func (t *teams) get(id db.TeamID) (team *db.Team, ok bool) {
	t.RLock()
	team, ok = t.m[id]
	t.RUnlock()
	return
}
