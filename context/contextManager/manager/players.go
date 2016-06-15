package manager

import (
	"encoding/json"
	"sync"

	"github.com/frrakn/treebeer/context/db"
	"github.com/frrakn/treebeer/context/position"
	ctxPb "github.com/frrakn/treebeer/context/proto"
	"github.com/juju/errors"
)

type players struct {
	m map[db.LcsID]*db.Player
	sync.RWMutex
}

func playerPbToDb(player *ctxPb.Player, id db.PlayerID) (*db.Player, error) {
	addlpos := make([]position.Position, len(player.Addlpos))
	for i, pos := range player.Addlpos {
		addlpos[i] = position.FromString(pos)
	}

	addlposJSON, err := json.Marshal(addlpos)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &db.Player{
		PlayerID: id,
		LcsID:    db.LcsID(player.Lcsid),
		RiotID:   db.RiotID(player.Riotid),
		Name:     player.Name,
		Position: position.FromString(player.Position),
		AddlPos:  string(addlposJSON),
	}, nil
}

func playerPbEqualsDb(update *ctxPb.Player, existing *db.Player) (bool, error) {
	addlpos := make([]position.Position, len(update.Addlpos))
	for i, pos := range update.Addlpos {
		addlpos[i] = position.FromString(pos)
	}

	addlposJSON, err := json.Marshal(addlpos)
	if err != nil {
		return false, errors.Trace(err)
	}

	return (int32(update.Lcsid) == int32(existing.LcsID) &&
		int32(update.Riotid) == int32(existing.RiotID) &&
		update.Name == existing.Name &&
		int32(update.Teamid) == int32(existing.TeamID) &&
		update.Position == existing.Position.String() &&
		string(addlposJSON) == existing.AddlPos), nil
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
