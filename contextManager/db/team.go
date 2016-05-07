package db

import (
	"fmt"

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

func CreateTeams(tx *sqlx.Tx, teams []*Team) error {
	for _, team := range teams {
		id, err := team.Create(tx)
		if err != nil {
			return errors.Annotate(err, fmt.Sprintf("Unable to create team with id %d", id))
		}
	}

	return nil
}

func UpdateTeams(tx *sqlx.Tx, teams []*Team) error {
	for _, team := range teams {
		err := team.Update(tx)
		if err != nil {
			return errors.Annotate(err, fmt.Sprintf("Unable to update team with id %d", team.TeamId))
		}
	}

	return nil
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
		return 0, err
	}

	id, err := res.LastInsertId()

	return TeamID(id), err
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

	return err
}
