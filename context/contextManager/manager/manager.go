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
	SEASONUPDATE_TAG   = "CtxMgr<SeasonUpdate>"
	LIVESTATUPDATE_TAG = "CtxMgr<LiveStatUpdate>"
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

func (s *Server) GetPlayer(ctx context.Context, player *ctxPb.Player) (*ctxPb.SavedPlayer, error) {
	p, ok := s.players.get(db.LcsID(player.Lcsid))
	if ok {
		return p.ToPB()
	}

	p = &db.Player{}
	teamID, err := s.teams.convertID(db.LcsID(player.Teamid))
	if err != nil {
		return nil, errors.Trace(err)
	}
	p.FromPB(player, 0)
	p.SetTeamID(teamID)
	err = db.EditTransact(s.SqlDB, LIVESTATUPDATE_TAG+" Players", func(tx *sqlx.Tx) error {
		_, err := p.Create(tx)
		if err != nil {
			return errors.Trace(err)
		}

		return db.UpdateVersion(tx, db.PlayerTable)
	})
	if err != nil {
		return nil, errors.Trace(err)
	}

	s.players.set(p.LcsID, p)
	return p.ToPB()
}

func (s *Server) GetTeam(ctx context.Context, team *ctxPb.Team) (*ctxPb.SavedTeam, error) {
	t, ok := s.teams.get(db.LcsID(team.Lcsid))
	if ok {
		return t.ToPB(), nil
	}

	t = &db.Team{}
	t.FromPB(team, 0)
	err := db.EditTransact(s.SqlDB, LIVESTATUPDATE_TAG+" Teams", func(tx *sqlx.Tx) error {
		_, err := t.Create(tx)
		if err != nil {
			return errors.Trace(err)
		}

		return db.UpdateVersion(tx, db.TeamTable)
	})
	if err != nil {
		return nil, errors.Trace(err)
	}

	s.teams.set(t.LcsID, t)
	return t.ToPB(), nil
}

func (s *Server) GetGame(ctx context.Context, game *ctxPb.Game) (*ctxPb.SavedGame, error) {
	g, ok := s.games.get(db.LcsID(game.Lcsid))
	if ok && g.EqualsPB(game) {
		return g.ToPB(), nil
	}

	g = &db.Game{}
	g.FromPB(game, 0)
	err := db.EditTransact(s.SqlDB, LIVESTATUPDATE_TAG+" Games", func(tx *sqlx.Tx) error {
		_, err := g.Create(tx)
		if err != nil {
			return errors.Trace(err)
		}

		return db.UpdateVersion(tx, db.GameTable)
	})
	if err != nil {
		return nil, errors.Trace(err)
	}

	s.games.set(g.LcsID, g)
	return g.ToPB(), nil
}

func (s *Server) GetStat(ctx context.Context, stat *ctxPb.Stat) (*ctxPb.SavedStat, error) {
	st, ok := s.stats.get(stat.Riotname)
	if ok {
		return st.ToPB(), nil
	}

	st = &db.Stat{}
	st.FromPB(stat, 0)
	err := db.EditTransact(s.SqlDB, LIVESTATUPDATE_TAG+" Stats", func(tx *sqlx.Tx) error {
		_, err := st.Create(tx)
		if err != nil {
			return errors.Trace(err)
		}

		return db.UpdateVersion(tx, db.StatTable)
	})
	if err != nil {
		return nil, errors.Trace(err)
	}

	s.stats.set(st.RiotName, st)
	return st.ToPB(), nil
}

func (s *Server) SeasonUpdate(ctx context.Context, updates *ctxPb.SeasonUpdates) (*ctxPb.Empty, error) {
	tDiff, pDiff, err := s.seasonUpdateDiff(updates)
	if err != nil {
		return empty, errors.Trace(err)
	}

	if len(tDiff.create) != 0 || len(tDiff.update) != 0 {
		// First, write teams to db
		err = db.EditTransact(s.SqlDB, SEASONUPDATE_TAG+" Teams", func(tx *sqlx.Tx) error {
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
			// TODO(frrakn): this should really be a batchupdate and the version update should ocur withing the db package...
			if len(tDiff.create) > 0 || len(tDiff.update) > 0 {
				err := db.UpdateVersion(tx, db.TeamTable)
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
		s.teams.batchUpdate(tDiff.create)
		s.teams.batchUpdate(tDiff.update)
	}

	if len(pDiff.create) != 0 || len(pDiff.update) != 0 {
		// Set player team IDs
		for _, p := range pDiff.create {
			p.TeamID, err = s.teams.convertID(p.teamLcsID)
			if err != nil {
				return nil, errors.Trace(err)
			}
		}
		for _, p := range pDiff.update {
			p.TeamID, err = s.teams.convertID(p.teamLcsID)
			if err != nil {
				return nil, errors.Trace(err)
			}
		}

		// Write players to db
		err = db.EditTransact(s.SqlDB, SEASONUPDATE_TAG+" Players", func(tx *sqlx.Tx) error {
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

			// TODO(frrakn): this should really be a batchupdate and the version update should ocur withing the db package...
			if len(pDiff.create) > 0 || len(pDiff.update) > 0 {
				err := db.UpdateVersion(tx, db.PlayerTable)
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

		// Update players in memory if db write was successful
		s.players.batchUpdateWithID(pDiff.create)
		s.players.batchUpdateWithID(pDiff.update)
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
		team, ok := s.teams.get(db.LcsID(t.Lcsid))
		dbTeam := new(db.Team)
		if !ok {
			dbTeam.FromPB(t, 0)
			tDiff.create = append(tDiff.create, dbTeam)
		} else if !team.EqualsPB(t) {
			dbTeam.FromPB(t, team.TeamID)
			tDiff.update = append(tDiff.update, dbTeam)
		}
	}

	for _, p := range ps {
		player, ok := s.players.get(db.LcsID(p.Lcsid))
		if !ok {
			dbPlayer := new(db.Player)
			err := dbPlayer.FromPB(p, 0)
			if err != nil {
				return nil, nil, errors.Trace(err)
			}
			pDiff.create = append(pDiff.create, &playerAndTeamID{
				teamLcsID: db.LcsID(p.Teamid),
				Player:    dbPlayer,
			})
		} else {
			equals, err := player.EqualsPB(p)
			if err != nil {
				return nil, nil, errors.Trace(err)
			}
			if !equals {
				dbPlayer := new(db.Player)
				err := dbPlayer.FromPB(p, player.PlayerID)
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
