package db

import "github.com/jmoiron/sqlx"

type PlayerID int32

type Player struct {
	PlayerId PlayerID
	LcsId    int32
	RiotId   int32
	Name     string
	TeamId   TeamID
}

func (p *Player) Create(tx *sqlx.Tx) (PlayerID, error) {
	res, err := tx.Exec(`
		INSERT INTO players
		VALUES (?, ?, ?, ?, ?)
		`,
		p.PlayerId,
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
