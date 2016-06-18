package poller

import (
	"sync"
	"time"

	"github.com/frrakn/treebeer/context/db"
)

type Poller struct {
	// TODO(fchen): this should really not be returning db.SeasonContext but something more along the lines of a protobuf... our services shouldn't be speaking in db layer language but rather protobuf. Should be a global change
	getUpdates func() (*db.SeasonContext, error)
	onClose    func()
	interval   *interval
	stop       chan struct{}
	Updates    chan *db.SeasonContext
	Errors     chan error
}

type interval struct {
	t time.Duration
	sync.RWMutex
}

const (
	DEFAULT_INTERVAL = time.Second * 10
)

func NewPoller(updates func() (*db.SeasonContext, error)) *Poller {
	return &Poller{
		getUpdates: updates,
		onClose: func() {
			return
		},
		interval: &interval{
			t: DEFAULT_INTERVAL,
		},
		stop:    make(chan struct{}),
		Updates: make(chan *db.SeasonContext),
		Errors:  make(chan error),
	}
}

func (p *Poller) SetInterval(t time.Duration) {
	p.interval.set(t)
}

func (p *Poller) SetClose(close func()) {
	p.onClose = close
}

func (p *Poller) Run() {
	go func() {
		for {
			interval := p.interval.get()
			time.Sleep(interval)

			updates, err := p.getUpdates()
			p.Errors <- err
			p.Updates <- updates
		}
	}()

	<-p.stop
}

func (p *Poller) Stop() {
	p.onClose()
	close(p.stop)
	close(p.Updates)
}

func (i *interval) set(t time.Duration) {
	i.Lock()
	i.t = t
	i.Unlock()
}

func (i *interval) get() time.Duration {
	i.RLock()
	t := i.t
	i.RUnlock()
	return t
}
