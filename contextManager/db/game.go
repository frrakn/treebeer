package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
)

type GameID int32

type Game struct {
	GameId      GameID
	LcsId       int32
	RiotGameId  string
	RiotMatchId string
	RedTeamId   TeamID
	BlueTeamId  TeamID
	GameStart   int64
	GameEnd     int64
}

func CreateGames(tx *sqlx.Tx, games []*Game) error {
	for _, game := range games {
		id, err := game.Create(tx)
		if err != nil {
			return errors.Annotate(err, fmt.Sprintf("Unable to create game with id %d", id))
		}
	}

	return nil
}

func UpdateGames(tx *sqlx.Tx, games []*Game) error {
	for _, game := range games {
		err := game.Update(tx)
		if err != nil {
			return errors.Annotate(err, fmt.Sprintf("Unable to update game with id %d", game.GameId))
		}
	}

	return nil
}

func (g *Game) Create(tx *sqlx.Tx) (GameID, error) {
	res, err := tx.Exec(`
		INSERT INTO games
		VALUES (NULL, ?, ?, ?, ?, ?, ?, ?)
		`,
		g.LcsId,
		g.RiotGameId,
		g.RiotMatchId,
		g.RedTeamId,
		g.BlueTeamId,
		g.GameStart,
		g.GameEnd,
	)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	return GameID(id), err
}

func (g *Game) Update(tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		UPDATE games
		SET
			lcsid = ?,
			riotgameid = ?,
			riotmatchid = ?,
			redteamid = ?,
			blueteamid = ?,
			gamestart = ?,
			gameend = ?
		WHERE gameid = ?
		`,
		g.LcsId,
		g.RiotGameId,
		g.RiotMatchId,
		g.RedTeamId,
		g.BlueTeamId,
		g.GameStart,
		g.GameEnd,
		g.GameId,
	)

	return err
}
