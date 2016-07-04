package aggregator

import (
	"sync"

	pb "github.com/frrakn/treebeer/liveData/liveData"
)

type Aggregator struct {
	m map[key]value
	sync.RWMutex
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		m: make(map[key]value),
	}
}

func (a *Aggregator) AddProto(p *pb.Stat) {
	key, value := fromPB(p)
	a.Lock()
	a.m[key] = value
	a.Unlock()
}

func (a *Aggregator) AddProtos(ps []*pb.Stat) {
	for _, p := range ps {
		a.AddProto(p)
	}
}

func (a *Aggregator) ToProtos() []*pb.Stat {
	protos := make([]*pb.Stat, 0)
	a.RLock()
	for k, v := range a.m {
		protos = append(protos, toPB(k, v))
	}
	a.RUnlock()
	return protos
}

func (a *Aggregator) Zero() {
	a.Lock()
	a.m = make(map[key]value)
	a.Unlock()
}
