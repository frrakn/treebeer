package server

import (
	"sync"

	"github.com/frrakn/treebeer/context/db"
)

type stats struct {
	m map[db.StatID]*db.Stat
	sync.RWMutex
}

func (s *stats) batchUpdate(stats []*db.Stat) {
	s.Lock()
	for _, stat := range stats {
		s.m[stat.StatID] = stat
	}
	s.Unlock()
}

func (s *stats) get(id db.StatID) (stat *db.Stat, ok bool) {
	s.RLock()
	stat, ok = s.m[id]
	s.RUnlock()
	return
}
