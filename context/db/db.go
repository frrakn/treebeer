package db

import (
	_ "github.com/go-sql-driver/mysql"
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

type Keyfiles struct {
	CaCert     string
	ClientCert string
	ClientKey  string
}

type Version struct {
	TableName string
	Version   int32
}

const (
	DB_STR_LEN      = 40
	INITIAL_VERSION = 1
)

func InitDB(dsn string) (*sqlx.DB, error) {
	sqldb, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, errors.Annotatef(err, "Unable to connect to database at %s", dsn)
	}

	return sqldb, nil
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

			season.Games, err = AllGames(tx)
			if err != nil {
				return errors.Trace(err)
			}

			season.Stats, err = AllStats(tx)
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
func EditTransact(db *sqlx.DB, tag string, action func(tx *sqlx.Tx) error) (err error) {
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

	return updateTimestamp(tx, tag)
}

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

	return
}

func updateTimestamp(tx *sqlx.Tx, tag string) error {

	if len(tag) > DB_STR_LEN {
		return errors.Errorf("tag is too long to insert into DB")
	}

	_, err := tx.Exec(`
		INSERT INTO updates (tag)
		VALUES (?)
	`, tag)

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

func UpdateVersion(tx *sqlx.Tx, table string) error {
	_, err := tx.Exec(`
		INSERT INTO versions (tablename, version)
		VALUES (?, ?)
		ON DUPLICATE KEY
		UPDATE version=version+1;
	`, table, INITIAL_VERSION)

	return errors.Trace(err)
}

func GetVersions(db *sqlx.DB) (map[string]int32, error) {
	var vs []*Version
	err := db.Select(&vs, `
		SELECT *
		FROM versions
	`)
	if err != nil {
		return nil, errors.Trace(err)
	}

	versions := make(map[string]int32)

	for _, v := range vs {
		versions[v.TableName] = v.Version
	}

	return versions, nil
}
