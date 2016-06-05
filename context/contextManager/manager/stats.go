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
