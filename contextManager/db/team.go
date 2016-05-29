package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
)

type TeamID int32

type Team struct {
	TeamId TeamID
	LcsId  int32
	RiotId int32
	Name   string
	Tag    string
}

func AllTeams(tx *sqlx.Tx) ([]*Team, error) {
	var teams []*Team

	err := tx.Select(&teams, `
		SELECT *
		FROM teams
	`)

	if err != nil {
		return nil, errors.Trace(err)
	}

	return teams, nil
}

func (t *Team) Create(tx *sqlx.Tx) (TeamID, error) {
	res, err := tx.Exec(`
		INSERT INTO teams
		VALUES (NULL, ?, ?, ?, ?)
		`,
		t.LcsId,
		t.RiotId,
		t.Name,
		t.Tag,
	)

	if err != nil {
		return 0, errors.Trace(err)
	}

	id, err := res.LastInsertId()

	return TeamID(id), errors.Trace(err)
}

func (t *Team) Update(tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		UPDATE teams
		SET
			lcsid = ?,
			riotid = ?,
			name = ?,
			tag = ?,
		WHERE teamid = ?
		`,
		t.LcsId,
		t.RiotId,
		t.Name,
		t.Tag,
		t.TeamId,
	)

	return errors.Trace(err)
}
