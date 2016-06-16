package manager

import (
	"sync"

	"github.com/frrakn/treebeer/context/db"
)

type teams struct {
	m map[db.LcsID]*db.Team
	sync.RWMutex
}

func (t *teams) batchUpdate(teams []*db.Team) {
	t.Lock()
	for _, team := range teams {
		t.m[team.LcsID] = team
	}
	t.Unlock()
}

func (t *teams) get(id db.LcsID) (team *db.Team, ok bool) {
	t.RLock()
	team, ok = t.m[id]
	t.RUnlock()
	return
}

func (t *teams) convertID(id db.LcsID) db.TeamID {
	t.RLock()
	teamID := t.m[id].TeamID
	t.RUnlock()
	return teamID
}
