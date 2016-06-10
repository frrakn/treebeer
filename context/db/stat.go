package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
)

type StatID int32

type Stat struct {
	StatID   StatID
	RiotName string
}

const (
	StatTable = "stats"
)

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
	if s.StatID != 0 {
		return 0, errors.Errorf("Stat has already been created!")
	}

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
	if err != nil {
		return 0, errors.Trace(err)
	}

	s.StatID = StatID(id)

	return s.StatID, nil
}

func (s *Stat) Update(tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		UPDATE stats
		SET
			riotname = ?
		WHERE statid = ?
		`,
		s.RiotName,
		s.StatID,
	)

	return errors.Trace(err)
}
