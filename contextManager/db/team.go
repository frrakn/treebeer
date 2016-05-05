package db

import "github.com/jmoiron/sqlx"

type TeamID int32

type Team struct {
	TeamId TeamID
	LcsId  int32
	RiotId int32
	Name   string
	Tag    string
}

func (t *Team) Create(tx *sqlx.Tx) (TeamID, error) {
	res, err := tx.Exec(`
		INSERT INTO teams
		VALUES (?, ?, ?, ?, ?)
		`,
		t.TeamId,
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
