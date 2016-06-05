package manager

import (
	"golang.org/x/net/context"

	"github.com/frrakn/treebeer/context/db"
	ctxPb "github.com/frrakn/treebeer/context/proto"
	"github.com/frrakn/treebeer/util/handle"
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
)

type Server struct {
	SqlDB   *sqlx.DB
	players *players
	teams   *teams
	games   *games
	stats   *stats
}

type teamDiff struct {
	create []*db.Team
	update []*db.Team
}

type playerDiff struct {
	create []*playerAndTeamID
	update []*playerAndTeamID
}

type playerAndTeamID struct {
	*db.Player
	teamLcsID db.LcsID
}

const (
	PROCESS_TAG = "CtxMgr"
)

var (
	empty = &ctxPb.Empty{}
)

func NewServer() *Server {
	return &Server{
		players: &players{
			m: make(map[db.LcsID]*db.Player),
		},
		teams: &teams{
			m: make(map[db.LcsID]*db.Team),
		},
		games: &games{
			m: make(map[db.LcsID]*db.Game),
		},
		stats: &stats{
			m: make(map[string]*db.Stat),
		},
	}

}

func (s *Server) Initialize(update *db.SeasonContext) {
	s.players.batchUpdate(update.Players)
	s.teams.batchUpdate(update.Teams)
	s.games.batchUpdate(update.Games)
	s.stats.batchUpdate(update.Stats)
}

func (s *Server) SeasonUpdate(ctx context.Context, updates *ctxPb.SeasonUpdates) (*ctxPb.Empty, error) {
	tDiff, pDiff, err := s.seasonUpdateDiff(updates)
	if err != nil {
		return empty, errors.Trace(err)
	}

	if len(tDiff.create) != 0 && len(tDiff.update) != 0 {
		// First, write teams to db
		err = db.EditTransact(s.SqlDB, PROCESS_TAG+" Teams", func(tx *sqlx.Tx) error {
			for _, ct := range tDiff.create {
				_, err := ct.Create(tx)
				if err != nil {
					return errors.Trace(err)
				}
			}
			for _, ut := range tDiff.update {
				err := ut.Update(tx)
				if err != nil {
					return errors.Trace(err)
				}
			}

			return nil
		})
		if err != nil {
			handle.Error(errors.Trace(err))
			return empty, errors.Trace(err)
		}

		// Update teams in memory if db write was successful
		s.teams.Lock()
		for _, ct := range tDiff.create {
			s.teams.m[ct.LcsID] = ct
		}
		for _, ut := range tDiff.update {
			s.teams.m[ut.LcsID] = ut
		}
		s.teams.Unlock()
	}

	if len(pDiff.create) != 0 && len(pDiff.update) != 0 {
		// Set player team IDs
		s.teams.RLock()
		for _, p := range pDiff.create {
			p.TeamID = s.teams.m[p.teamLcsID].TeamID
		}
		for _, p := range pDiff.update {
			p.TeamID = s.teams.m[p.teamLcsID].TeamID
		}
		s.teams.RUnlock()

		// Write players to db
		err = db.EditTransact(s.SqlDB, PROCESS_TAG+" Players", func(tx *sqlx.Tx) error {
			for _, cp := range pDiff.create {
				_, err := cp.Create(tx)
				if err != nil {
					return errors.Trace(err)
				}
			}
			for _, up := range pDiff.update {
				err := up.Update(tx)
				if err != nil {
					return errors.Trace(err)
				}
			}

			return nil
		})
		if err != nil {
			handle.Error(errors.Trace(err))
			return empty, errors.Trace(err)
		}

		s.players.Lock()
		// Update players in memory if db write was successful
		for _, cp := range pDiff.create {
			s.players.m[cp.LcsID] = cp.Player
		}
		for _, up := range pDiff.update {
			s.players.m[up.LcsID] = up.Player
		}
		s.players.Unlock()
	}

	return empty, errors.Trace(err)
}

func (s *Server) seasonUpdateDiff(updates *ctxPb.SeasonUpdates) (tDiff *teamDiff, pDiff *playerDiff, err error) {
	ts := updates.Teams
	ps := updates.Players

	tDiff = &teamDiff{
		create: []*db.Team{},
		update: []*db.Team{},
	}

	pDiff = &playerDiff{
		create: []*playerAndTeamID{},
		update: []*playerAndTeamID{},
	}

	for _, t := range ts {
		s.teams.RLock()
		team, ok := s.teams.m[db.LcsID(t.Lcsid)]
		s.teams.RUnlock()
		if !ok {
			tDiff.create = append(tDiff.create, teamPbToDb(t, 0))
		} else if !teamPbEqualsDb(t, team) {
			tDiff.update = append(tDiff.update, teamPbToDb(t, team.TeamID))
		}
	}

	for _, p := range ps {
		s.players.RLock()
		player, ok := s.players.m[db.LcsID(p.Lcsid)]
		s.players.RUnlock()
		if !ok {
			dbPlayer, err := playerPbToDb(p, 0)
			if err != nil {
				return nil, nil, errors.Trace(err)
			}
			pDiff.create = append(pDiff.create, &playerAndTeamID{
				teamLcsID: db.LcsID(p.Teamid),
				Player:    dbPlayer,
			})
		} else {
			equals, err := playerPbEqualsDb(p, player)
			if err != nil {
				return nil, nil, errors.Trace(err)
			}
			if !equals {
				dbPlayer, err := playerPbToDb(p, player.PlayerID)
				if err != nil {
					return nil, nil, errors.Trace(err)
				}
				pDiff.update = append(pDiff.update, &playerAndTeamID{
					teamLcsID: db.LcsID(p.Teamid),
					Player:    dbPlayer,
				})
			}
		}
	}

	return
}
