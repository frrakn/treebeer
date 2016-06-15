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

func (t *teams) all() []*db.Team {
	ts := []*db.Team{}
	t.RLock()
	for _, team := range t.m {
		ts = append(ts, team)
	}
	t.RUnlock()
	return ts
}
