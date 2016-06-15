package server

import (
	"fmt"
	"net"

	"github.com/frrakn/treebeer/context/db"
	ctxPb "github.com/frrakn/treebeer/context/proto"
	"github.com/frrakn/treebeer/util/handle"
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Server struct {
	poller  *DBPoller
	players *players
	teams   *teams
	games   *games
	stats   *stats
	stop    chan struct{}
}

func NewServer(sqlDB *sqlx.DB) *Server {
	s := &Server{
		poller: NewDBPoller(sqlDB),
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

	return s
}

func (s *Server) Run(port string) {
	go s.poller.Run()
	go s.serveRpc(port)
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

func (s *Server) serveRpc(port string) {
	l, err := net.Listen("tcp", port)
	if err != nil {
		handle.Fatal(errors.Annotate(err, fmt.Sprintf("Unable to get listener on port %d", port)))
	}

	rpcserv := grpc.NewServer()
	ctxPb.RegisterSeasonContextServer(rpcserv, s)
	rpcserv.Serve(l)
}

func (s *Server) GetContext(_ context.Context, _ *ctxPb.Empty) (*ctxPb.FullContext, error) {
	players := s.players.all()
	teams := s.teams.all()
	games := s.games.all()
	stats := s.stats.all()

	fullCtx := &ctxPb.FullContext{
		Players: make([]*ctxPb.SavedPlayer, len(players)),
		Teams:   make([]*ctxPb.SavedTeam, len(teams)),
		Games:   make([]*ctxPb.SavedGame, len(games)),
		Stats:   make([]*ctxPb.SavedStat, len(stats)),
	}

	var err error
	for i, player := range players {
		fullCtx.Players[i], err = player.ToPB()
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	for i, team := range teams {
		fullCtx.Teams[i] = team.ToPB()
	}

	for i, game := range games {
		fullCtx.Games[i] = game.ToPB()
	}

	for i, stat := range stats {
		fullCtx.Stats[i] = stat.ToPB()
	}

	return fullCtx, nil
}

func (s *Server) Stop() {
	close(s.stop)
}
