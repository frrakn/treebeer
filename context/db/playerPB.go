package db

import (
	"encoding/json"

	"github.com/frrakn/treebeer/context/position"
	ctxPb "github.com/frrakn/treebeer/context/proto"
	"github.com/juju/errors"
)

func (p *Player) ToPB() (*ctxPb.SavedPlayer, error) {
	var ap []position.Position
	err := json.Unmarshal([]byte(p.AddlPos), &ap)
	if err != nil {
		return nil, errors.Trace(err)
	}
	addlpos := make([]string, len(ap))
	for i, pos := range ap {
		addlpos[i] = pos.String()
	}

	player := &ctxPb.SavedPlayer{
		Player: &ctxPb.Player{
			Lcsid:    int32(p.LcsID),
			Riotid:   int32(p.RiotID),
			Name:     p.Name,
			Teamid:   int32(p.TeamID),
			Position: p.Position.String(),
			Addlpos:  addlpos,
		},
		Playerid: int32(p.PlayerID),
	}
	return player, nil
}

func (p *Player) FromPB(player *ctxPb.Player, id PlayerID) error {
	addlpos := make([]position.Position, len(player.Addlpos))
	for i, pos := range player.Addlpos {
		addlpos[i] = position.FromString(pos)
	}

	addlposJSON, err := json.Marshal(addlpos)
	if err != nil {
		return errors.Trace(err)
	}

	p.PlayerID = id
	p.LcsID = LcsID(player.Lcsid)
	p.RiotID = RiotID(player.Riotid)
	p.Name = player.Name
	// TeamID foreign key omitted for internal translation
	p.Position = position.FromString(player.Position)
	p.AddlPos = string(addlposJSON)

	return nil
}

func (p *Player) FromSavedPB(player *ctxPb.SavedPlayer) error {
	addlpos := make([]position.Position, len(player.Player.Addlpos))
	for i, pos := range player.Player.Addlpos {
		addlpos[i] = position.FromString(pos)
	}

	addlposJSON, err := json.Marshal(addlpos)
	if err != nil {
		return errors.Trace(err)
	}

	p.PlayerID = PlayerID(player.Playerid)
	p.LcsID = LcsID(player.Player.Lcsid)
	p.RiotID = RiotID(player.Player.Riotid)
	p.Name = player.Player.Name
	p.TeamID = TeamID(player.Player.Teamid)
	p.Position = position.FromString(player.Player.Position)
	p.AddlPos = string(addlposJSON)

	return nil
}

func (p *Player) EqualsPB(player *ctxPb.Player) (bool, error) {
	addlpos := make([]position.Position, len(player.Addlpos))
	for i, pos := range player.Addlpos {
		addlpos[i] = position.FromString(pos)
	}

	addlposJSON, err := json.Marshal(addlpos)
	if err != nil {
		return false, errors.Trace(err)
	}

	return (int32(player.Lcsid) == int32(p.LcsID) &&
		int32(player.Riotid) == int32(p.RiotID) &&
		player.Name == p.Name &&
		int32(player.Teamid) == int32(p.TeamID) &&
		player.Position == p.Position.String() &&
		string(addlposJSON) == p.AddlPos), nil
}

func (p *Player) SetTeamID(team TeamID) {
	p.TeamID = team
}
