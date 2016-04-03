/*
==========
AGGREGATOR
==========
Provides data manipulation service for incoming map[string]string

*/

package aggregator

type Aggregator struct {
	sum map[string]int
}

func NewAggregator() *Aggregator {
	a := new(Aggregator)
	a.sum = make(map[string]int)
}

func (a *Aggregator) Add(newAgg map[string]int) {
	for key, value := range newAgg {
		a.sum[key] += value
	}
}

func (a *Aggregator) GetSum() map[string]int {
	return a.sum
}
