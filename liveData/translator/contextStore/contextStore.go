package contextStore

import (
	"github.com/frrakn/treebeer/context/db"
	"github.com/frrakn/treebeer/context/poller"
	pb "github.com/frrakn/treebeer/context/proto"
	"github.com/juju/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type ContextStore struct {
	players *players
	teams   *teams
	games   *games
	stats   *stats

	poller    *poller.RPCPoller
	ctxClient pb.LiveStatUpdateClient

	conn *grpc.ClientConn

	config *Configuration
	ready  chan struct{}
	stop   chan struct{}
	Errors chan error
}

type Configuration struct {
	RPCPollAddr string
	CtxMgrAddr  string
}

func New(cfg *Configuration) *ContextStore {
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
	// TODO(frrakn): This makes process hang if ready is never received... stuck listening for ContextCache
	<-c.ready
}

func (c *ContextStore) Run() {
	var err error
	c.poller, err = poller.NewRPCPoller(c.config.RPCPollAddr)
	if err != nil {
		c.Errors <- errors.Trace(err)
		c.Stop()
		return
	}

	go c.poller.Run()

	c.conn, err = grpc.Dial(c.config.CtxMgrAddr, grpc.WithInsecure())
	if err != nil {
		c.Errors <- errors.Trace(err)
		c.Stop()
		return
	}

	c.ctxClient = pb.NewLiveStatUpdateClient(c.conn)

	firstUpdate := true
	for {
		select {
		case <-c.stop:
			return
		case season := <-c.poller.Updates:
			c.players.batchUpdate(season.Players)
			c.teams.batchUpdate(season.Teams)
			c.games.batchUpdate(season.Games)
			c.stats.batchUpdate(season.Stats)
			if firstUpdate {
				close(c.ready)
				firstUpdate = false
			}
		case err = <-c.poller.Errors:
			if err != nil {
				c.Errors <- errors.Trace(err)
			}
		}
	}
}

func (c *ContextStore) Stop() {
	c.conn.Close()
	c.poller.Stop()
	close(c.stop)
}

func (c *ContextStore) ConvertPlayer(id db.RiotID) (db.PlayerID, error) {
	player, ok := c.players.get(id)
	if ok {
		return player.PlayerID, nil
	}

	return 0, errors.Errorf("Unable to translate player with ID: %d", id)
}

func (c *ContextStore) GetTeamForPlayer(id db.RiotID) (db.TeamID, error) {
	player, ok := c.players.get(id)
	if !ok {
		return 0, errors.Errorf("Unable to translate player with ID: %d", id)
	}

	return player.TeamID, nil
}

func (c *ContextStore) ConvertTeam(id db.RiotID) (db.TeamID, error) {
	team, ok := c.teams.get(id)
	if ok {
		return team.TeamID, nil
	}

	return 0, errors.Errorf("Unable to translate team with ID: %d", id)
}

func (c *ContextStore) ConvertGame(g *pb.Game) (db.GameID, error) {
	game, ok := c.games.get(db.LcsID(g.Lcsid))
	if ok && game.EqualsPB(g) {
		return game.GameID, nil
	}

	sgame, err := c.ctxClient.GetGame(context.Background(), g)
	if err != nil {
		return 0, errors.Errorf("Unable to translate game: %s", g)
	}

	return db.GameID(sgame.Gameid), nil
}

func (c *ContextStore) ConvertStat(s *pb.Stat) (db.StatID, error) {
	stat, ok := c.stats.get(s.Riotname)
	if ok {
		return stat.StatID, nil
	}

	sStat, err := c.ctxClient.GetStat(context.Background(), s)
	if err != nil {
		return 0, errors.Errorf("Unable to translate stat: %s", s)
	}

	return db.StatID(sStat.Statid), nil
}
