package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
)

type PlayerID int32

type Player struct {
	PlayerId PlayerID
	LcsId    int32
	RiotId   int32
	Name     string
	TeamId   TeamID
}

func AllPlayers(tx *sqlx.Tx) ([]*Player, error) {
	var players []*Player

	err := tx.Select(&players, `
		SELECT *
		FROM players
	`)

	if err != nil {
		return nil, errors.Trace(err)
	}

	return players, nil
}

func (p *Player) Create(tx *sqlx.Tx) (PlayerID, error) {
	res, err := tx.Exec(`
		INSERT INTO players
		VALUES (NULL, ?, ?, ?, ?)
		`,
		p.LcsId,
		p.RiotId,
		p.Name,
		p.TeamId,
	)

	if err != nil {
		return 0, errors.Trace(err)
	}

	id, err := res.LastInsertId()

	return PlayerID(id), errors.Trace(err)
}

func (p *Player) Update(tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		UPDATE players
		SET
			lcsid = ?,
			riotid = ?,
			name = ?,
			teamid = ?,
		WHERE playerid = ?
		`,
		p.LcsId,
		p.RiotId,
		p.Name,
		p.TeamId,
		p.PlayerId,
	)

	return errors.Trace(err)
}
