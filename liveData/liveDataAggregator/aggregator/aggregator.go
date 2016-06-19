/*
==========
AGGREGATOR
==========
Provides data manipulation service for incoming map[string]string

*/

package aggregator

type Aggregator struct {
	sum  map[string]int
	sema chan struct{}
}

func NewAggregator() *Aggregator {
	a := new(Aggregator)
	a.sum = make(map[string]int)
	a.sema = make(chan struct{}, 1)
	a.sema <- struct{}{}
	return a
}

func (a *Aggregator) Add(newAgg map[string]int) {
	<-a.sema
	defer func() {
		a.sema <- struct{}{}
	}()
	for key, value := range newAgg {
		a.sum[key] += value
	}
}

func (a *Aggregator) GetSum() map[string]int {
	<-a.sema
	defer func() {
		a.sema <- struct{}{}
	}()
	return a.sum
}

func (a *Aggregator) Zero() {
	<-a.sema
	defer func() {
		a.sema <- struct{}{}
	}()
	a.sum = make(map[string]int)
}
