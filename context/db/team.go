package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
)

type TeamID int32

type Team struct {
	TeamID TeamID
	LcsID  LcsID
	RiotID RiotID
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
	if t.TeamID != 0 {
		return 0, errors.Errorf("Team has already been created!")
	}

	res, err := tx.Exec(`
		INSERT INTO teams
		VALUES (NULL, ?, ?, ?, ?)
		`,
		t.LcsID,
		t.RiotID,
		t.Name,
		t.Tag,
	)

	if err != nil {
		return 0, errors.Trace(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.Trace(err)
	}

	t.TeamID = TeamID(id)

	return t.TeamID, nil
}

func (t *Team) Update(tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		UPDATE teams
		SET
			lcsid = ?,
			riotid = ?,
			name = ?,
			tag = ?
		WHERE teamid = ?
		`,
		t.LcsID,
		t.RiotID,
		t.Name,
		t.Tag,
		t.TeamID,
	)

	return errors.Trace(err)
}
