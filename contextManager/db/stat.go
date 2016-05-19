package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
)

type StatID int32

type Stat struct {
	StatId   StatID
	RiotName string
}

func CreateStats(tx *sqlx.Tx, stats []*Stat) error {
	for _, stat := range stats {
		id, err := stat.Create(tx)
		if err != nil {
			return errors.Annotate(err, fmt.Sprintf("Unable to create stat with id %d", id))
		}
	}

	return nil
}

func UpdateStats(tx *sqlx.Tx, stats []*Stat) error {
	for _, stat := range stats {
		err := stat.Update(tx)
		if err != nil {
			return errors.Annotate(err, fmt.Sprintf("Unable to update stat with id %d", stat.StatId))
		}
	}

	return nil
}

func (s *Stat) Create(tx *sqlx.Tx) (StatID, error) {
	res, err := tx.Exec(`
		INSERT INTO stats
		VALUES (NULL, ?)
		`,
		s.RiotName,
	)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	return StatID(id), err
}

func (s *Stat) Update(tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		UPDATE stats
		SET
			riotname = ?,
		WHERE statid = ?
		`,
		s.RiotName,
		s.StatId,
	)

	return err
}
