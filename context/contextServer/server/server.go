package server

import (
	"net/http"

	"github.com/frrakn/treebeer/context/db"
	"github.com/frrakn/treebeer/context/poller"
	"github.com/frrakn/treebeer/util/handle"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	Router  *httprouter.Router
	poller  *poller.DBPoller
	players *players
	teams   *teams
	games   *games
	stats   *stats
	stop    chan struct{}
}

func NewServer(sqlDB *sqlx.DB) *Server {
	s := &Server{
		Router: httprouter.New(),
		poller: poller.NewDBPoller(sqlDB),
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
		stop: make(chan struct{}),
	}

	s.Router.GET("/player/:id", getPlayer(s.players))
	s.Router.GET("/team/:id", getTeam(s.teams))
	s.Router.GET("/game/:id", getGame(s.games))
	s.Router.GET("/stat/:id", getStat(s.stats))

	return s
}

func (s *Server) Run(port string) {
	go s.poller.Run()
	go http.ListenAndServe(port, s.Router)
	defer s.poller.Stop()
	var (
		season *db.SeasonContext
		err    error
	)
	for {
		select {
		case <-s.stop:
			return
		case season = <-s.poller.Updates:
			s.players.batchUpdate(season.Players)
			s.teams.batchUpdate(season.Teams)
			s.games.batchUpdate(season.Games)
			s.stats.batchUpdate(season.Stats)
		case err = <-s.poller.Errors:
			if err != nil {
				handle.Error(err)
			}
		}
	}
}

func (s *Server) Stop() {
	close(s.stop)
}
