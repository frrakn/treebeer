package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
)

// Transact wraps transactional db functions that begins and commits / rollbacks transactions
func Transact(db *sqlx.DB, action func(tx *sqlx.Tx) error) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		err = errors.Trace(err)
		return
	}

	defer func() {
		if pnc := recover(); pnc != nil {
			switch pnc := pnc.(type) {
			case error:
				err = errors.Trace(pnc)
			default:
				err = errors.Errorf("%s", pnc)
			}
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()
	return action(tx)
}
