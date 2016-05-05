package db

import "github.com/jmoiron/sqlx"

type GameID int32

type Game struct {
	GameId      GameID
	LcsId       int32
	RiotGameId  string
	RiotMatchId string
	RedTeamId   TeamID
	BlueTeamId  TeamID
	GameStart   int64
	GameEnd     int64
}

func (g *Game) Create(tx *sqlx.Tx) (GameID, error) {
	res, err := tx.Exec(`
		INSERT INTO games
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`,
		g.GameId,
		g.LcsId,
		g.RiotGameId,
		g.RiotMatchId,
		g.RedTeamId,
		g.BlueTeamId,
		g.GameStart,
		g.GameEnd,
	)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	return GameID(id), err
}

func (g *Game) Update(tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		UPDATE games
		SET
			lcsid = ?,
			riotgameid = ?,
			riotmatchid = ?,
			redteamid = ?,
			blueteamid = ?,
			gamestart = ?,
			gameend = ?
		WHERE gameid = ?
		`,
		g.LcsId,
		g.RiotGameId,
		g.RiotMatchId,
		g.RedTeamId,
		g.BlueTeamId,
		g.GameStart,
		g.GameEnd,
		g.GameId,
	)

	return err
}
