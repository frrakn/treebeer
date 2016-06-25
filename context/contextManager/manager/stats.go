package manager

import (
	"sync"

	"github.com/frrakn/treebeer/context/db"
)

type stats struct {
	m map[string]*db.Stat
	sync.RWMutex
}

func (s *stats) batchUpdate(stats []*db.Stat) {
	s.Lock()
	for _, stat := range stats {
		s.m[stat.RiotName] = stat
	}
	s.Unlock()
}

func (s *stats) get(id string) (stat *db.Stat, ok bool) {
	s.RLock()
	stat, ok = s.m[id]
	s.RUnlock()
	return
}

func (s *stats) set(name string, stat *db.Stat) {
	s.RLock()
	s.m[name] = stat
	s.RUnlock()
}
