package db

import ctxPb "github.com/frrakn/treebeer/context/proto"

func (t *Team) ToPB() *ctxPb.SavedTeam {
	return &ctxPb.SavedTeam{
		Team: &ctxPb.Team{
			Lcsid:  int32(t.LcsID),
			Riotid: int32(t.RiotID),
			Name:   t.Name,
			Tag:    t.Tag,
		},
		Teamid: int32(t.TeamID),
	}
}

func (t *Team) FromPB(team *ctxPb.Team, id TeamID) {
	t.TeamID = id
	t.LcsID = LcsID(team.Lcsid)
	t.RiotID = RiotID(team.Riotid)
	t.Name = team.Name
	t.Tag = team.Tag
}

func (t *Team) EqualsPB(team *ctxPb.Team) bool {
	return int32(team.Lcsid) == int32(t.LcsID) &&
		int32(team.Riotid) == int32(t.RiotID) &&
		team.Name == t.Name &&
		team.Tag == t.Tag
}
