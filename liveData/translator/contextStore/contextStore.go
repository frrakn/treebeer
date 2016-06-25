package contextStore

import (
	"github.com/frrakn/treebeer/context/db"
	"github.com/juju/errors"
	"google.golang.org/grpc"
)

type ContextStore struct {
	players *players
	teams   *teams
	games   *games
	stats   *stats

	poller    *RPCPoller
	ctxClient *pb.LiveStatUpdateClient

	config Configuration
	ready  chan struct{}
	stop   chan struct{}
	Errors chan error
}

type Configuration struct {
	pollerAddr string
	ctxMgrAddr string
}

func NewContextStore(cfg Configuration) (*ContextStore, error) {
	return &ContextStore{
		players: newPlayers(),
		teams:   newTeams(),
		games:   newGames(),
		stats:   newStats(),
		config:  cfg,
		ready:   make(chan struct{}),
		stop:    make(chan struct{}),
		Errors:  make(chan error),
	}
}

func (c *ContextStore) Start() {
	go c.Run()
	<-ready
}

func (c *ContextStore) Run() {
	poller, err := NewRPCPoller(c.config.pollerAddr)
	if err != nil {
		c.Errors <- errors.Trace(err)
		c.Stop()
		return
	}

	poller.Run()

	conn, err := grpc.Dial(c.config.ctxMgrAddr)
	defer conn.Close()
	if err != nil {
		c.Errors <- errors.Trace(err)
		c.Stop()
		return
	}

	c.ctxClient = pb.NewLiveStatUpdateClient(conn)

	firstUpdate := true
	for {
		select {
		case <-c.stop:
			return
		case season = <-c.poller.Updates:
			c.players.batchUpdate(season.Players)
			c.teams.batchUpdate(season.Teams)
			c.games.batchUpdate(season.Games)
			c.stats.batchUpdate(season.Stats)
			if firstUpdate {
				close(c.ready)
			}
		case err = <-c.poller.Errors:
			if err != nil {
				c.Errors <- errors.Trace(err)
			}
		}
	}
}

func (c *ContextStore) Stop() {
	c.poller.Stop()
	close(c.stop)
}

func (c *ContextStore) convertPlayer(id db.LcsID) db.PlayerID {
	player, ok := p.players.get(id)

}

func (c *ContextStore) convertTeam(id db.LcsID) db.TeamID {

}

func (c *ContextStore) convertGame(id db.LcsID) db.GameID {

}

func (c *ContextStore) convertStat(name string) db.StatID {

}
