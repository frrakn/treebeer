package server

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/frrakn/treebeer/context/db"
	pb "github.com/frrakn/treebeer/context/proto"
	"github.com/frrakn/treebeer/liveData/translator/contextStore"
	"github.com/frrakn/treebeer/liveData/translator/ws"
	"github.com/frrakn/treebeer/liveData/translator/ws/schema"
	"github.com/juju/errors"
)

type Server struct {
	config *Configuration

	ctxStore *contextStore.ContextStore
	listener *ws.Listener

	stop   chan struct{}
	Errors chan error
}

type Configuration struct {
	Port string
}

type Stat struct {
	game   db.GameID
	team   db.TeamID
	player db.PlayerID
	stat   db.StatID
	val    json.RawMessage
	t      int32
}

var (
	MAX_TIME = time.Unix(1<<63-1, 0)
)

func NewServer(cfg *Configuration, ctxStore *contextStore.ContextStore, listener *ws.Listener) *Server {
	server := &Server{
		config:   cfg,
		ctxStore: ctxStore,
		listener: listener,
		stop:     make(chan struct{}),
		Errors:   make(chan error),
	}

	return server
}

func (s *Server) Start() {
	go s.run()
}

func (s *Server) run() {
	go s.handlers()
	s.ctxStore.Start()
	s.listener.Start()
}

func (s *Server) Stop() {
	s.ctxStore.Stop()
	s.listener.Stop()
	close(s.stop)
}

func (s *Server) handlers() {
	for {
		select {
		case err := <-s.ctxStore.Errors:
			s.Errors <- errors.Trace(err)
		case err := <-s.listener.Errors:
			s.Errors <- errors.Trace(err)
		case stats := <-s.listener.Stats:
			err := s.handleStats(stats)
			if err != nil {
				s.Errors <- errors.Trace(err)
			}
		case <-s.stop:
			return
		}
	}
}

func (s *Server) handleStats(stats map[string]*schema.Game) error {
	for gameid, game := range stats {
		lcsid, err := strconv.Atoi(gameid)
		if err != nil {
			return errors.Trace(err)
		}

		bluePlayer, ok := game.PlayerStats["1"]
		if !ok {
			return errors.Errorf("Unexpected schema change, blue player not found")
		}
		bluePlayerID, err := strconv.Atoi(bluePlayer.PlayerID)
		if err != nil {
			return errors.Trace(err)
		}

		blueTeamID, err := s.ctxStore.GetTeamForPlayer(db.RiotID(bluePlayerID))
		if err != nil {
			return errors.Trace(err)
		}

		redPlayer, ok := game.PlayerStats["6"]
		if !ok {
			return errors.Errorf("Unexpected schema change, red player not found")
		}
		redPlayerID, err := strconv.Atoi(redPlayer.PlayerID)
		if err != nil {
			return errors.Trace(err)
		}

		redTeamID, err := s.ctxStore.GetTeamForPlayer(db.RiotID(redPlayerID))
		if err != nil {
			return errors.Trace(err)
		}

		g := &pb.Game{
			Lcsid:       int32(lcsid),
			Riotgameid:  game.GameID,
			Riotmatchid: game.MatchID,
			Redteamid:   int32(redTeamID),
			Blueteamid:  int32(blueTeamID),
			Gamestart:   time.Now().Unix(),
			Gameend:     MAX_TIME.Unix(),
		}

		_, err = s.ctxStore.ConvertGame(g)
	}

	return nil
}
