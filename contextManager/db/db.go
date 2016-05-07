package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
)

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

	return err
}

func UnsafeFkCheck(tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		SET FOREIGN_KEY_CHECKS = 0
	`)

	return err
}

func SafeFkCheck(tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		SET FOREIGN_KEY_CHECKS = 1
	`)

	return err
}
