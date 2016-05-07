package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/frrakn/treebeer/contextManager/db"
	pb "github.com/frrakn/treebeer/contextManager/proto"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"

	"golang.org/x/net/context"
)

type configuration struct {
	DB   string
	Port string
}

type server struct {
	sqldb *sqlx.DB
}

type batchUpdate struct {
	teamsCreate   []*db.Team
	teamsUpdate   []*db.Team
	playersCreate []*db.Player
	playersUpdate []*db.Player
	gamesCreate   []*db.Game
	gamesUpdate   []*db.Game
}

func main() {
	flag.Parse()

	c := loadConfig()
	db := initDB(c.DB)
	serveRpc(c.Port, db)
}

func handleFatal(err error) {
	glog.Fatal("%s", err)
}

func handleError(err error) {
	glog.Error("%s", err)
}

func errorString(err error) string {
	return fmt.Sprintf("%s", err)
}

func loadConfig() configuration {
	var c configuration
	if len(os.Args) < 2 {
		handleFatal(errors.Errorf("No config file specified"))
	}

	fileLoc := os.Args[1]
	configBytes, err := ioutil.ReadFile(fileLoc)
	if err != nil {
		handleFatal(errors.Annotate(err, fmt.Sprintf("Error reading from file location %s", fileLoc)))
	}

	err = json.Unmarshal(configBytes, &c)
	if err != nil {
		handleFatal(errors.Annotate(err, fmt.Sprintf("Error parsing config json from %s into config struct", fileLoc)))
	}

	return c
}

func initDB(dsn string) *sqlx.DB {
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		handleFatal(errors.Annotate(err, fmt.Sprintf("Unable to connect to database at %s", dsn)))
	}

	return db
}

func serveRpc(port string, db *sqlx.DB) {
	l, err := net.Listen("tcp", port)
	if err != nil {
		handleFatal(errors.Annotate(err, fmt.Sprintf("Unable to get listener on port %d", port)))
	}

	rpcserv := grpc.NewServer()
	pb.RegisterContextUpdateServer(rpcserv, &server{db})
	rpcserv.Serve(l)
}

func teamPbToDb(team *pb.Team) *db.Team {
	return &db.Team{db.TeamID(team.Teamid), team.Lcsid, team.Riotid, team.Name, team.Tag}
}

func teamsPbToDb(teams []*pb.Team) []*db.Team {
	ts := make([]*db.Team, len(teams))
	for i, team := range teams {
		ts[i] = teamPbToDb(team)
	}
	return ts
}

func playerPbToDb(player *pb.Player) *db.Player {
	return &db.Player{db.PlayerID(player.Playerid), player.Lcsid, player.Riotid, player.Name, db.TeamID(player.Teamid)}
}

func playersPbToDb(players []*pb.Player) []*db.Player {
	ps := make([]*db.Player, len(players))
	for i, player := range players {
		ps[i] = playerPbToDb(player)
	}
	return ps
}

func gamePbToDb(game *pb.Game) *db.Game {
	return &db.Game{db.GameID(game.Gameid), game.Lcsid, game.Riotgameid, game.Riotmatchid, db.TeamID(game.Redteamid), db.TeamID(game.Blueteamid), game.Gamestart, game.Gameend}
}

func gamesPbToDb(games []*pb.Game) []*db.Game {
	gs := make([]*db.Game, len(games))
	for i, game := range games {
		gs[i] = gamePbToDb(game)
	}
	return gs
}

func batchPbToDb(batch *pb.BatchUpdates) *batchUpdate {
	return &batchUpdate{
		teamsPbToDb(batch.TeamsCreate.Teams),
		teamsPbToDb(batch.TeamsUpdate.Teams),
		playersPbToDb(batch.PlayersCreate.Players),
		playersPbToDb(batch.PlayersUpdate.Players),
		gamesPbToDb(batch.GamesCreate.Games),
		gamesPbToDb(batch.GamesUpdate.Games),
	}
}

func (s *server) CreateTeam(ctx context.Context, team *pb.Team) (*pb.Result, error) {
	t := teamPbToDb(team)
	err := db.Transact(s.sqldb, func(tx *sqlx.Tx) error {
		_, err := t.Create(tx)
		return err
	})

	if err != nil {
		return &pb.Result{pb.Result_FAIL, errorString(errors.Annotate(err, fmt.Sprintf("Error creating new team with struct %+v", team)))}, err
	}

	return &pb.Result{pb.Result_SUCCESS, ""}, nil
}

func (s *server) UpdateTeam(ctx context.Context, team *pb.Team) (*pb.Result, error) {
	t := teamPbToDb(team)
	err := db.Transact(s.sqldb, func(tx *sqlx.Tx) error {
		return t.Update(tx)
	})

	if err != nil {
		return &pb.Result{pb.Result_FAIL, errorString(errors.Annotate(err, fmt.Sprintf("Error updating team with struct %+v", team)))}, err
	}

	return &pb.Result{pb.Result_SUCCESS, ""}, nil
}

func (s *server) CreatePlayer(ctx context.Context, player *pb.Player) (*pb.Result, error) {
	p := playerPbToDb(player)
	err := db.Transact(s.sqldb, func(tx *sqlx.Tx) error {
		_, err := p.Create(tx)
		return err
	})

	if err != nil {
		return &pb.Result{pb.Result_FAIL, errorString(errors.Annotate(err, fmt.Sprintf("Error creating new player with struct %+v", player)))}, err
	}

	return &pb.Result{pb.Result_SUCCESS, ""}, nil
}

func (s *server) UpdatePlayer(ctx context.Context, player *pb.Player) (*pb.Result, error) {
	p := playerPbToDb(player)
	err := db.Transact(s.sqldb, func(tx *sqlx.Tx) error {
		return p.Update(tx)
	})

	if err != nil {
		return &pb.Result{pb.Result_FAIL, errorString(errors.Annotate(err, fmt.Sprintf("Error updating player with struct %+v", player)))}, err
	}

	return &pb.Result{pb.Result_SUCCESS, ""}, nil
}

func (s *server) CreateGame(ctx context.Context, game *pb.Game) (*pb.Result, error) {
	g := gamePbToDb(game)
	err := db.Transact(s.sqldb, func(tx *sqlx.Tx) error {
		_, err := g.Create(tx)
		return err
	})

	if err != nil {
		return &pb.Result{pb.Result_FAIL, errorString(errors.Annotate(err, fmt.Sprintf("Error creating new game with struct %+v", game)))}, err
	}

	return &pb.Result{pb.Result_SUCCESS, ""}, nil
}

func (s *server) UpdateGame(ctx context.Context, game *pb.Game) (*pb.Result, error) {
	g := gamePbToDb(game)
	err := db.Transact(s.sqldb, func(tx *sqlx.Tx) error {
		return g.Update(tx)
	})

	if err != nil {
		return &pb.Result{pb.Result_FAIL, errorString(errors.Annotate(err, fmt.Sprintf("Error updating game with struct %+v", game)))}, err
	}

	return &pb.Result{pb.Result_SUCCESS, ""}, nil
}

func (s *server) BatchUpdate(ctx context.Context, update *pb.BatchUpdates) (*pb.Result, error) {
	batch := batchPbToDb(update)
	err := db.Transact(s.sqldb, func(tx *sqlx.Tx) error {
		err := db.CreateTeams(tx, batch.teamsCreate)
		if err != nil {
			return err
		}

		err = db.UpdateTeams(tx, batch.teamsUpdate)
		if err != nil {
			return err
		}

		err = db.CreatePlayers(tx, batch.playersCreate)
		if err != nil {
			return err
		}

		err = db.UpdatePlayers(tx, batch.playersUpdate)
		if err != nil {
			return err
		}

		err = db.CreateGames(tx, batch.gamesCreate)
		if err != nil {
			return err
		}

		return db.UpdateGames(tx, batch.gamesUpdate)
	})

	if err != nil {
		return &pb.Result{pb.Result_FAIL, errorString(errors.Annotate(err, fmt.Sprintf("Error updating batch with struct %+v", update)))}, err
	}

	return &pb.Result{pb.Result_SUCCESS, ""}, nil
}
