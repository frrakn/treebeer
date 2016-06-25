package contextStore

import (
	"sync"

	"github.com/frrakn/treebeer/context/db"
)

type stats struct {
	m map[string]*db.Stat
	sync.RWMutex
}

func newStats() *stats {
	return &stats{
		m: make(map[string]*db.Stat),
	}
}

func (s *stats) batchUpdate(stats []*db.Stat) {
	s.Lock()
	for _, stat := range stats {
		s.m[stat.RiotName] = stat
	}
	s.Unlock()
}

func (s *stats) get(name string) (stat *db.Stat, ok bool) {
	s.RLock()
	stat, ok = s.m[name]
	s.RUnlock()
	return
}

func (s *stats) set(name string, stat *db.Stat) {
	s.RLock()
	s.m[name] = stat
	s.RUnlock()
}
