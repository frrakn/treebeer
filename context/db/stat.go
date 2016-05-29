package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
)

type StatID int32

type Stat struct {
	StatId   StatID
	RiotName string
}

func AllStats(tx *sqlx.Tx) ([]*Stat, error) {
	var stats []*Stat

	err := tx.Select(&stats, `
		SELECT *
		FROM stats
	`)

	if err != nil {
		return nil, errors.Trace(err)
	}

	return stats, nil
}

func (s *Stat) Create(tx *sqlx.Tx) (StatID, error) {
	res, err := tx.Exec(`
		INSERT INTO stats
		VALUES (NULL, ?)
		`,
		s.RiotName,
	)

	if err != nil {
		return 0, errors.Trace(err)
	}

	id, err := res.LastInsertId()

	return StatID(id), errors.Trace(err)
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

	return errors.Trace(err)
}
