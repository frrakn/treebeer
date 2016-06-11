package server

import (
	"github.com/frrakn/treebeer/context/db"
	"github.com/frrakn/treebeer/context/poller"
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"
)

type DBPoller struct {
	*poller.Poller
	SqlDB    *sqlx.DB
	versions map[string]int32
}

func NewDBPoller(sqlDB *sqlx.DB) *DBPoller {
	p := &DBPoller{
		SqlDB:    sqlDB,
		versions: make(map[string]int32),
	}

	p.Poller = poller.NewPoller(
		func() (*db.SeasonContext, error) {
			season := &db.SeasonContext{}
			currVersions, err := db.GetVersions(p.SqlDB)
			if err != nil {
				return nil, errors.Trace(err)
			}

			currPlayerVersion, ok := currVersions[db.PlayerTable]
			if ok {
				playerVersion, ok := p.versions[db.PlayerTable]
				if !ok || currPlayerVersion != playerVersion {
					err := db.Transact(p.SqlDB, func(tx *sqlx.Tx) error {
						var err error
						season.Players, err = db.AllPlayers(tx)
						return errors.Trace(err)
					})
					if err != nil {
						return nil, errors.Trace(err)
					}
				}
			}

			currTeamVersion, ok := currVersions[db.TeamTable]
			if ok {
				teamVersion, ok := p.versions[db.TeamTable]
				if !ok || currTeamVersion != teamVersion {
					err := db.Transact(p.SqlDB, func(tx *sqlx.Tx) error {
						var err error
						season.Teams, err = db.AllTeams(tx)
						return errors.Trace(err)
					})
					if err != nil {
						return nil, errors.Trace(err)
					}
				}
			}

			currGameVersion, ok := currVersions[db.GameTable]
			if ok {
				gameVersion, ok := p.versions[db.GameTable]
				if !ok || currGameVersion != gameVersion {
					err := db.Transact(p.SqlDB, func(tx *sqlx.Tx) error {
						var err error
						season.Games, err = db.AllGames(tx)
						return errors.Trace(err)
					})
					if err != nil {
						return nil, errors.Trace(err)
					}
				}
			}

			currStatVersion, ok := currVersions[db.StatTable]
			if ok {
				statVersion, ok := p.versions[db.StatTable]
				if !ok || currStatVersion != statVersion {
					err := db.Transact(p.SqlDB, func(tx *sqlx.Tx) error {
						var err error
						season.Stats, err = db.AllStats(tx)
						return errors.Trace(err)
					})
					if err != nil {
						return nil, errors.Trace(err)
					}
				}
			}

			return season, nil
		},
	)

	return p
}
