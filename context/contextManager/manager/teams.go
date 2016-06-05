package manager

import (
	"sync"

	"github.com/frrakn/treebeer/context/db"
	ctxPb "github.com/frrakn/treebeer/context/proto"
)

type teams struct {
	m map[db.LcsID]*db.Team
	sync.RWMutex
}

func teamPbToDb(team *ctxPb.Team, id db.TeamID) *db.Team {
	return &db.Team{
		TeamID: id,
		LcsID:  db.LcsID(team.Lcsid),
		RiotID: db.RiotID(team.Riotid),
		Name:   team.Name,
		Tag:    team.Tag,
	}
}

func teamPbEqualsDb(update *ctxPb.Team, existing *db.Team) bool {
	return int32(update.Lcsid) == int32(existing.LcsID) &&
		int32(update.Riotid) == int32(existing.RiotID) &&
		update.Name == existing.Name &&
		update.Tag == existing.Tag
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
