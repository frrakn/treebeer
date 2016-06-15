package db

import (
	"encoding/json"

	"github.com/frrakn/treebeer/context/position"
	ctxPb "github.com/frrakn/treebeer/context/proto"
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
)

type PlayerID int32

type Player struct {
	PlayerID PlayerID
	LcsID    LcsID
	RiotID   RiotID
	Name     string
	TeamID   TeamID
	Position position.Position
	AddlPos  string
}

const (
	PlayerTable = "players"
)

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

func (p *Player) ToPB() (*ctxPb.SavedPlayer, error) {
	var ap []position.Position
	err := json.Unmarshal([]byte(p.AddlPos), &ap)
	if err != nil {
		return nil, errors.Trace(err)
	}
	addlpos := make([]string, len(ap))
	for i, pos := range ap {
		addlpos[i] = pos.String()
	}

	player := &ctxPb.SavedPlayer{
		Player: &ctxPb.Player{
			Lcsid:    int32(p.LcsID),
			Riotid:   int32(p.RiotID),
			Name:     p.Name,
			Teamid:   int32(p.TeamID),
			Position: p.Position.String(),
			Addlpos:  addlpos,
		},
		Playerid: int32(p.PlayerID),
	}
	return player, nil
}

func (p *Player) Create(tx *sqlx.Tx) (PlayerID, error) {
	if p.PlayerID != 0 {
		return 0, errors.Errorf("Player has already been created!")
	}

	res, err := tx.Exec(`
		INSERT INTO players
		VALUES (NULL, ?, ?, ?, ?, ?, ?)
		`,
		p.LcsID,
		p.RiotID,
		p.Name,
		p.TeamID,
		p.Position,
		p.AddlPos,
	)

	if err != nil {
		return 0, errors.Trace(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.Trace(err)
	}

	p.PlayerID = PlayerID(id)

	return p.PlayerID, nil
}

func (p *Player) Update(tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		UPDATE players
		SET
			lcsid = ?,
			riotid = ?,
			name = ?,
			teamid = ?,
			position = ?,
			addlpos = ?
		WHERE playerid = ?
		`,
		p.LcsID,
		p.RiotID,
		p.Name,
		p.TeamID,
		p.Position,
		p.AddlPos,
		p.PlayerID,
	)

	return errors.Trace(err)
}
