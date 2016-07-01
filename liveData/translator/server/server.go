package server

import (
	"bufio"
	"encoding/json"
	"net"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/frrakn/treebeer/context/db"
	ctxPb "github.com/frrakn/treebeer/context/proto"
	ldPb "github.com/frrakn/treebeer/liveData/liveData"
	"github.com/frrakn/treebeer/liveData/translator/contextStore"
	"github.com/frrakn/treebeer/liveData/translator/ws"
	"github.com/frrakn/treebeer/liveData/translator/ws/schema"
	"github.com/juju/errors"
)

type Server struct {
	config *Configuration

	ctxStore     *contextStore.ContextStore
	riotListener *ws.Listener

	servListener net.Listener
	connections  *connections

	stop   chan struct{}
	Errors chan error
}

type Configuration struct {
	Port string
}

type connections struct {
	sync.Mutex
	c []net.Conn
}

var (
	MAX_TIME = time.Unix(1<<63-1, 0)
)

func NewServer(cfg *Configuration, ctxStore *contextStore.ContextStore, listener *ws.Listener) *Server {
	server := &Server{
		config:       cfg,
		ctxStore:     ctxStore,
		riotListener: listener,
		connections: &connections{
			c: make([]net.Conn, 0),
		},
		stop:   make(chan struct{}),
		Errors: make(chan error),
	}

	return server
}

func (s *Server) Start() {
	go s.run()
}

func (s *Server) run() {
	go s.handlers()
	s.ctxStore.Start()
	s.riotListener.Start()

	var err error
	s.servListener, err = net.Listen("tcp", s.config.Port)
	if err != nil {
		s.Stop()
	}

	go s.acceptLoop()
}

func (s *Server) Stop() {
	s.ctxStore.Stop()
	s.riotListener.Stop()
	if s.servListener != nil {
		s.servListener.Close()
	}
	s.connections.closeAll()
	close(s.stop)
}

func (s *Server) handlers() {
	for {
		select {
		case err := <-s.ctxStore.Errors:
			s.Errors <- errors.Trace(err)
		case err := <-s.riotListener.Errors:
			s.Errors <- errors.Trace(err)
		case stats := <-s.riotListener.Stats:
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

		g := &ctxPb.Game{
			Lcsid:       int32(lcsid),
			Riotgameid:  game.GameID,
			Riotmatchid: game.MatchID,
			Redteamid:   int32(redTeamID),
			Blueteamid:  int32(blueTeamID),
			Gamestart:   time.Now().Unix(),
			Gameend:     MAX_TIME.Unix(),
		}

		gameID, err := s.ctxStore.ConvertGame(g)
		timestamp := game.T

		for _, team := range game.TeamStats {
			var teamID db.TeamID
			if team.Color == "blue" {
				teamID = blueTeamID
			} else {
				teamID = redTeamID
			}
			teamStruct := reflect.Indirect(reflect.ValueOf(team))
			for i := 0; i < teamStruct.NumField(); i++ {
				statName := teamStruct.Type().Field(i).Name
				statVal := teamStruct.Field(i).Interface()
				statBytes, err := json.Marshal(statVal)
				if err != nil {
					return errors.Trace(err)
				}
				statID, err := s.ctxStore.ConvertStat(
					&ctxPb.Stat{
						Riotname: statName,
					})
				if err != nil {
					return errors.Trace(err)
				}
				proto := ldPb.Stat{
					Playerid:  0,
					Teamid:    int32(teamID),
					Gameid:    int32(gameID),
					Statid:    int32(statID),
					Jsonvalue: string(statBytes),
					Timestamp: timestamp,
				}
				protoBytes, err := json.Marshal(proto)
				if err != nil {
					return errors.Trace(err)
				}

				s.broadcast(protoBytes)
			}
		}
		for _, player := range game.PlayerStats {
			playerID, err := strconv.ParseInt(player.PlayerID, 10, 32)
			if err != nil {
				return errors.Trace(err)
			}
			playerStruct := reflect.Indirect(reflect.ValueOf(player))
			for i := 0; i < playerStruct.NumField(); i++ {
				statName := playerStruct.Type().Field(i).Name
				statVal := playerStruct.Field(i).Interface()
				statBytes, err := json.Marshal(statVal)
				if err != nil {
					return errors.Trace(err)
				}
				statID, err := s.ctxStore.ConvertStat(
					&ctxPb.Stat{
						Riotname: statName,
					})
				if err != nil {
					return errors.Trace(err)
				}
				proto := ldPb.Stat{
					Playerid:  int32(playerID),
					Teamid:    0,
					Gameid:    int32(gameID),
					Statid:    int32(statID),
					Jsonvalue: string(statBytes),
					Timestamp: timestamp,
				}
				protoBytes, err := json.Marshal(proto)
				if err != nil {
					return errors.Trace(err)
				}

				s.broadcast(protoBytes)
			}
		}
	}

	return nil
}

func (s *Server) acceptLoop() {
	conn, err := s.servListener.Accept()
	if err != nil {
		s.Errors <- errors.Trace(err)
	}

	s.connections.add(conn)
}

func (s *Server) broadcast(msg []byte) {
	errs := s.connections.broadcast(msg)
	for _, err := range errs {
		s.Errors <- err
	}
}

func (c *connections) add(conn net.Conn) {
	c.Lock()
	c.c = append(c.c, conn)
	c.Unlock()
}

func (c *connections) broadcast(msg []byte) []error {
	errs := make([]error, 0)
	for i, conn := range c.c {
		writer := bufio.NewWriter(conn)

		_, err := writer.Write(msg)
		if err != nil {
			errs = append(errs, errors.Trace(err))
			c.closeAndRemove(i)
		}
	}

	return errs
}

func (c *connections) closeAndRemove(i int) {
	c.Lock()
	c.c[i].Close()
	c.c[i] = c.c[len(c.c)-1]
	c.c[len(c.c)-1] = nil
	c.c = c.c[:len(c.c)-1]
	c.Unlock()
}

func (c *connections) closeAll() {
	for _, conn := range c.c {
		conn.Close()
	}
}
