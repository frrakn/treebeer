package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
)

type LcsID int32
type RiotID int32

type SeasonContext struct {
	Players []*Player
	Teams   []*Team
	Stats   []*Stat
	Games   []*Game
}

func GetSeasonContext(db *sqlx.DB) (*SeasonContext, error) {
	season := &SeasonContext{}

	err := Transact(
		db,
		func(tx *sqlx.Tx) error {
			var err error

			season.Players, err = AllPlayers(tx)
			if err != nil {
				return errors.Trace(err)
			}

			season.Teams, err = AllTeams(tx)
			if err != nil {
				return errors.Trace(err)
			}

			return nil
		},
	)

	if err != nil {
		return nil, errors.Trace(err)
	}

	return season, nil
}

// Transact wraps transactional db functions that begins and commits / rollbacks transactions
// Also auto-updates versioning
func Transact(db *sqlx.DB, action func(tx *sqlx.Tx) error) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		err = errors.Trace(err)
		return
	}

	err = action(tx)
	if err != nil {
		err = errors.Trace(err)
		return
	}

	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = errors.Trace(p)
			default:
				err = errors.Errorf("%s", p)
			}
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	return updateTimestamp(tx)
}

func updateTimestamp(tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		INSERT INTO updates
		VALUES ()
	`)

	return errors.Trace(err)
}

func UnsafeFkCheck(tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		SET FOREIGN_KEY_CHECKS = 0
	`)

	return errors.Trace(err)
}

func SafeFkCheck(tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		SET FOREIGN_KEY_CHECKS = 1
	`)

	return errors.Trace(err)
}
