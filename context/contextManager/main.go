package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/frrakn/treebeer/context/db"
	"github.com/frrakn/treebeer/context/position"
	ctxPb "github.com/frrakn/treebeer/context/proto"
	"github.com/frrakn/treebeer/util/config"
	"github.com/frrakn/treebeer/util/handle"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
)

type configuration struct {
	DB       string
	Port     string
	Keyfiles keyfiles
}

type keyfiles struct {
	CaCert     string
	ClientCert string
	ClientKey  string
}

type server struct {
	sqldb   *sqlx.DB
	players map[db.LcsID]*db.Player
	teams   map[db.LcsID]*db.Team
	games   map[db.LcsID]*db.Game
	stats   map[string]*db.Stat
}

var (
	empty     = &ctxPb.Empty{}
	conf      configuration
	ctxServer *server
)

func main() {
	serveRpc(conf.Port)
}

func init() {
	flag.Parse()

	ctxServer := &server{}

	err := config.LoadConfig(&conf)
	if err != nil {
		handle.Fatal(errors.Annotate(err, "Failed to load configuration"))
	}
	ctxServer.sqldb = initDB(conf.DB, conf.Keyfiles)

	season, err := db.GetSeasonContext(ctxServer.sqldb)
	if err != nil {
		handle.Fatal(errors.Annotate(err, "Failed to load season data from DB"))
	}

	ctxServer.players = make(map[db.LcsID]*db.Player)
	for _, p := range season.Players {
		ctxServer.players[p.LcsID] = p
	}

	ctxServer.teams = make(map[db.LcsID]*db.Team)
	for _, t := range season.Teams {
		ctxServer.teams[t.LcsID] = t
	}

	ctxServer.games = make(map[db.LcsID]*db.Game)
	for _, g := range season.Games {
		ctxServer.games[g.LcsID] = g
	}

	ctxServer.stats = make(map[string]*db.Stat)
	for _, s := range season.Stats {
		ctxServer.stats[s.RiotName] = s
	}
}

func initDB(dsn string, keys keyfiles) *sqlx.DB {
	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile(keys.CaCert)
	if err != nil {
		handle.Fatal(errors.Annotate(err, fmt.Sprintf("Unable to access database credentials at %s", keys.CaCert)))
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		handle.Fatal(errors.Annotate(err, "Unabe to append PEM."))
	}

	clientCert := make([]tls.Certificate, 0, 1)
	certs, err := tls.LoadX509KeyPair(keys.ClientCert, keys.ClientKey)
	if err != nil {
		handle.Fatal(errors.Annotate(err, fmt.Sprintf("Unable to access database credentials at %s and %s", keys.ClientCert, keys.ClientKey)))
	}
	clientCert = append(clientCert, certs)

	mysql.RegisterTLSConfig("treebeer", &tls.Config{
		RootCAs:            rootCertPool,
		Certificates:       clientCert,
		InsecureSkipVerify: true,
	})

	sqldb, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		handle.Fatal(errors.Annotate(err, fmt.Sprintf("Unable to connect to database at %s", dsn)))
	}

	return sqldb
}

func serveRpc(port string) {
	l, err := net.Listen("tcp", port)
	if err != nil {
		handle.Fatal(errors.Annotate(err, fmt.Sprintf("Unable to get listener on port %d", port)))
	}

	rpcserv := grpc.NewServer()
	ctxPb.RegisterSeasonUpdateServer(rpcserv, ctxServer)
	rpcserv.Serve(l)
}

func (s *server) SeasonUpdate(ctx context.Context, updates *ctxPb.SeasonUpdates) (*ctxPb.Empty, error) {
	createTeams, updateTeams, createPlayers, updatePlayers, err := s.seasonUpdateDiff(updates)
	if err != nil {
		return empty, errors.Trace(err)
	}

	// First, write teams to db
	err = db.Transact(s.sqldb, func(tx *sqlx.Tx) error {
		for _, ct := range createTeams {
			_, err := ct.Create(tx)
			if err != nil {
				return errors.Trace(err)
			}
		}
		for _, ut := range updateTeams {
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
	for _, ct := range createTeams {
		s.teams[ct.LcsID] = ct
	}
	for _, ut := range updateTeams {
		s.teams[ut.LcsID] = ut
	}

	// Set player team IDs
	for id, p := range createPlayers {
		p.TeamID = s.teams[id].TeamID
	}

	// Write players to db
	err = db.Transact(s.sqldb, func(tx *sqlx.Tx) error {
		for _, cp := range createPlayers {
			_, err := cp.Create(tx)
			if err != nil {
				return errors.Trace(err)
			}
		}
		for _, up := range updatePlayers {
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

	// Update players in memory if db write was successful
	for _, cp := range createPlayers {
		s.players[cp.LcsID] = cp
	}
	for _, up := range updatePlayers {
		s.players[up.LcsID] = up
	}

	return empty, errors.Trace(err)
}

func (s *server) seasonUpdateDiff(updates *ctxPb.SeasonUpdates) (createTeams []*db.Team, updateTeams []*db.Team, createPlayers map[db.LcsID]*db.Player, updatePlayers []*db.Player, err error) {
	ts := updates.Teams
	ps := updates.Players

	createTeams = []*db.Team{}
	updateTeams = []*db.Team{}
	createPlayers = make(map[db.LcsID]*db.Player)
	updatePlayers = []*db.Player{}

	for _, t := range ts {
		team, ok := s.teams[db.LcsID(t.Lcsid)]
		if !ok {
			createTeams = append(createTeams, teamPbToDb(t, 0))
		} else if !teamPbEqualsDb(t, team) {
			updateTeams = append(updateTeams, teamPbToDb(t, team.TeamID))
		}
	}

	for _, p := range ps {
		player, ok := s.players[db.LcsID(p.Lcsid)]
		if !ok {
			dbPlayer, err := playerPbToDb(p, 0)
			if err != nil {
				return nil, nil, nil, nil, err
			}
			createPlayers[db.LcsID(p.Lcsid)] = dbPlayer
		} else {
			equals, err := playerPbEqualsDb(p, player)
			if err != nil {
				return nil, nil, nil, nil, err
			}
			if !equals {
				dbPlayer, err := playerPbToDb(p, player.PlayerID)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				updatePlayers = append(updatePlayers, dbPlayer)
			}
		}
	}

	return
}

func teamPbToDb(team *ctxPb.Team, id db.TeamID) *db.Team {
	return &db.Team{
		TeamID: id,
		LcsID:  db.LcsID(team.Lcsid),
		RiotID: db.RiotID(team.Riotid),
		Name:   team.Name,
		Tag:    team.Tag,
	}
}

func teamPbEqualsDb(update *ctxPb.Team, existing *db.Team) bool {
	return int32(update.Lcsid) == int32(existing.LcsID) &&
		int32(update.Riotid) == int32(existing.RiotID) &&
		update.Name == existing.Name &&
		update.Tag == existing.Tag
}

func playerPbToDb(player *ctxPb.Player, id db.PlayerID) (*db.Player, error) {
	addlpos, err := json.Marshal(player.Addlpos)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &db.Player{
		PlayerID: id,
		LcsID:    db.LcsID(player.Lcsid),
		RiotID:   db.RiotID(player.Riotid),
		Name:     player.Name,
		Position: position.FromString(player.Position),
		AddlPos:  string(addlpos),
	}, nil
}

func playerPbEqualsDb(update *ctxPb.Player, existing *db.Player) (bool, error) {
	addlpos, err := json.Marshal(update.Addlpos)
	if err != nil {
		return false, errors.Trace(err)
	}
	return (int32(update.Lcsid) == int32(existing.LcsID) &&
		int32(update.Riotid) == int32(existing.RiotID) &&
		update.Name == existing.Name &&
		int32(update.Teamid) == int32(existing.TeamID) &&
		update.Position == existing.Position.String() &&
		string(addlpos) == existing.AddlPos), nil
}
