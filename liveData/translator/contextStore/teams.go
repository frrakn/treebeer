package contextStore

import (
	"sync"

	"github.com/frrakn/treebeer/context/db"
)

type teams struct {
	m map[db.RiotID]*db.Team
	sync.RWMutex
}

func newTeams() *teams {
	return &teams{
		m: make(map[db.RiotID]*db.Team),
	}
}

func (t *teams) batchUpdate(teams []*db.Team) {
	t.Lock()
	for _, team := range teams {
		t.m[team.RiotID] = team
	}
	t.Unlock()
}

func (t *teams) get(id db.RiotID) (team *db.Team, ok bool) {
	t.RLock()
	team, ok = t.m[id]
	t.RUnlock()
	return
}

func (t *teams) set(id db.RiotID, team *db.Team) {
	t.RLock()
	t.m[id] = team
	t.RUnlock()
}
