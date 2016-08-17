package schema

import (
	ctxPb "github.com/frrakn/treebeer/context/proto"
	"github.com/juju/errors"
)

type RiotSeason struct {
	SeasonName  string
	SeasonSplit string
	ProTeams    []RiotTeam
	ProPlayers  []RiotPlayer
}

type RiotTeam struct {
	Id        int32
	RiotId    int32
	Name      string
	ShortName string
}

type RiotPlayer struct {
	Id        int32
	RiotId    int32
	Name      string
	PhotoUrl  string
	ProTeamId int32
	Positions []string
}

func (rs *RiotSeason) ToSeasonUpdates() (*ctxPb.SeasonUpdates, error) {
	updates := &ctxPb.SeasonUpdates{
		Teams:   make([]*ctxPb.Team, len(rs.ProTeams)),
		Players: make([]*ctxPb.Player, len(rs.ProPlayers)),
	}

	for i, t := range rs.ProTeams {
		updates.Teams[i] = t.ToCtxTeam()
	}

	var err error
	for i, p := range rs.ProPlayers {
		updates.Players[i], err = p.ToCtxPlayer()
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	return updates, nil
}

func (rt *RiotTeam) ToCtxTeam() *ctxPb.Team {
	return &ctxPb.Team{
		Lcsid:  rt.Id,
		Riotid: rt.RiotId,
		Name:   rt.Name,
		Tag:    rt.ShortName,
	}
}

func (rp *RiotPlayer) ToCtxPlayer() (*ctxPb.Player, error) {
	if len(rp.Positions) == 0 {
		return nil, errors.Errorf("No positions parsed for player %s", rp.Name)
	}

	return &ctxPb.Player{
		Lcsid:    rp.Id,
		Riotid:   rp.RiotId,
		Name:     rp.Name,
		Photourl: rp.PhotoUrl,
		Teamid:   rp.ProTeamId,
		Position: rp.Positions[0],
		Addlpos:  rp.Positions[1:],
	}, nil
}
