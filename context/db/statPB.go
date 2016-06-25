package db

import ctxPb "github.com/frrakn/treebeer/context/proto"

func (s *Stat) ToPB() *ctxPb.SavedStat {
	return &ctxPb.SavedStat{
		Stat: &ctxPb.Stat{
			Riotname: s.RiotName,
		},
		Statid: int32(s.StatID),
	}
}

func (s *Stat) FromPB(stat *ctxPb.Stat, id StatID) {
	s.StatID = id
	s.RiotName = stat.Riotname
}

func (s *Stat) FromSavedPB(stat *ctxPb.SavedStat) {
	s.StatID = StatID(stat.Statid)
	s.RiotName = stat.Stat.Riotname
}
