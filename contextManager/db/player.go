package db

import (
	"fmt"

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

func CreatePlayers(tx *sqlx.Tx, players []*Player) error {
	for _, player := range players {
		id, err := player.Create(tx)
		if err != nil {
			return errors.Annotate(err, fmt.Sprintf("Unable to create player with id %d", id))
		}
	}

	return nil
}

func UpdatePlayers(tx *sqlx.Tx, players []*Player) error {
	for _, player := range players {
		err := player.Update(tx)
		if err != nil {
			return errors.Annotate(err, fmt.Sprintf("Unable to update player with id %d", player.PlayerId))
		}
	}

	return nil
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
		return 0, err
	}

	id, err := res.LastInsertId()

	return PlayerID(id), err
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

	return err
}
