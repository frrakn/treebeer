package server

import (
	"github.com/frrakn/treebeer/context/db"
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	SqlDB    *sqlx.DB
	Router   *httprouter.Router
	versions map[string]int32
	players  *players
	teams    *teams
	games    *games
	stats    *stats
}

func NewServer(sqlDB *sqlx.DB) *Server {
	s := &Server{
		SqlDB:    sqlDB,
		Router:   httprouter.New(),
		versions: make(map[string]int32),
		players: &players{
			m: make(map[db.PlayerID]*db.Player),
		},
		teams: &teams{
			m: make(map[db.TeamID]*db.Team),
		},
		games: &games{
			m: make(map[db.GameID]*db.Game),
		},
		stats: &stats{
			m: make(map[db.StatID]*db.Stat),
		},
	}

	s.Router.GET("/player/:id", getPlayer(s.players))
	s.Router.GET("/team/:id", getTeam(s.teams))
	s.Router.GET("/game/:id", getGame(s.games))
	s.Router.GET("/stat/:id", getStat(s.stats))

	return s
}

func (s *Server) Update() error {
	currVersions, err := db.GetVersions(s.SqlDB)
	if err != nil {
		return errors.Trace(err)
	}

	currPlayerVersion, ok := currVersions[db.PlayerTable]
	if ok {
		playerVersion, ok := s.versions[db.PlayerTable]
		if !ok || currPlayerVersion != playerVersion {
			err := s.updatePlayers()
			if err != nil {
				return errors.Trace(err)
			}
		}
	}

	currTeamVersion, ok := currVersions[db.TeamTable]
	if ok {
		teamVersion, ok := s.versions[db.TeamTable]
		if !ok || currTeamVersion != teamVersion {
			err := s.updateTeams()
			if err != nil {
				return errors.Trace(err)
			}
		}
	}

	currGameVersion, ok := currVersions[db.GameTable]
	if ok {
		gameVersion, ok := s.versions[db.GameTable]
		if !ok || currGameVersion != gameVersion {
			err := s.updateGames()
			if err != nil {
				return errors.Trace(err)
			}
		}
	}

	currStatVersion, ok := currVersions[db.StatTable]
	if ok {
		statVersion, ok := s.versions[db.StatTable]
		if !ok || currStatVersion != statVersion {
			err := s.updateStats()
			if err != nil {
				return errors.Trace(err)
			}
		}
	}

	return nil
}

func (s *Server) updatePlayers() error {
	var players []*db.Player
	err := db.Transact(s.SqlDB, func(tx *sqlx.Tx) error {
		var err error
		players, err = db.AllPlayers(tx)
		return errors.Trace(err)
	})

	s.players.batchUpdate(players)

	return errors.Trace(err)
}

func (s *Server) updateTeams() error {
	var teams []*db.Team
	err := db.Transact(s.SqlDB, func(tx *sqlx.Tx) error {
		var err error
		teams, err = db.AllTeams(tx)
		return errors.Trace(err)
	})

	s.teams.batchUpdate(teams)

	return errors.Trace(err)
}

func (s *Server) updateGames() error {
	var games []*db.Game
	err := db.Transact(s.SqlDB, func(tx *sqlx.Tx) error {
		var err error
		games, err = db.AllGames(tx)
		return errors.Trace(err)
	})

	s.games.batchUpdate(games)

	return errors.Trace(err)
}

func (s *Server) updateStats() error {
	var stats []*db.Stat
	err := db.Transact(s.SqlDB, func(tx *sqlx.Tx) error {
		var err error
		stats, err = db.AllStats(tx)
		return errors.Trace(err)
	})

	s.stats.batchUpdate(stats)

	return errors.Trace(err)
}
