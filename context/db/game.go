package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
)

type GameID int32

type Game struct {
	GameID      GameID
	LcsID       LcsID
	RiotGameID  string
	RiotMatchID string
	RedTeamID   TeamID
	BlueTeamID  TeamID
	GameStart   int64
	GameEnd     int64
}

func AllGames(tx *sqlx.Tx) ([]*Game, error) {
	var games []*Game

	err := tx.Select(&games, `
		SELECT *
		FROM games
	`)

	if err != nil {
		return nil, errors.Trace(err)
	}

	return games, nil
}

func (g *Game) Create(tx *sqlx.Tx) (GameID, error) {
	if g.GameID != 0 {
		return 0, errors.Errorf("Game has already been created!")
	}

	res, err := tx.Exec(`
		INSERT INTO games
		VALUES (NULL, ?, ?, ?, ?, ?, ?, ?)
		`,
		g.LcsID,
		g.RiotGameID,
		g.RiotMatchID,
		g.RedTeamID,
		g.BlueTeamID,
		g.GameStart,
		g.GameEnd,
	)

	if err != nil {
		return 0, errors.Trace(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.Trace(err)
	}

	g.GameID = GameID(id)

	return g.GameID, nil
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
		g.LcsID,
		g.RiotGameID,
		g.RiotMatchID,
		g.RedTeamID,
		g.BlueTeamID,
		g.GameStart,
		g.GameEnd,
		g.GameID,
	)

	return errors.Trace(err)
}
