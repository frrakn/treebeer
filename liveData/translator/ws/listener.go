package ws

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/frrakn/treebeer/liveData/translator/ws/schema"
	"github.com/gorilla/websocket"
	"github.com/juju/errors"
)

type Listener struct {
	address string
	path    string
	opts    map[string]string

	Errors chan error
	Stats  chan *schema.LiveStats
	stop   chan struct{}
}

type Configuration struct {
	Address string
	Path    string
	Opts    map[string]string
}

func NewListener(cfg *Configuration) *Listener {
	return &Listener{
		address: cfg.Address,
		path:    cfg.Path,
		opts:    cfg.Opts,
		Errors:  make(chan error),
		Stats:   make(chan *schema.LiveStats),
		stop:    make(chan struct{}),
	}
}

func (l *Listener) Start() {
	go l.Run()
}

func (l *Listener) Run() {
	header := http.Header{}
	for k, v := range l.opts {
		header.Add(k, v)
	}

	u := url.URL{Scheme: "ws", Host: l.address, Path: l.path}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		l.Errors <- errors.Trace(err)
	}
	defer conn.Close()

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				l.Errors <- errors.Trace(err)
				return
			}
			var liveStats schema.LiveStats
			err = json.Unmarshal(message, &liveStats)
			if err != nil {
				l.Errors <- errors.Trace(err)
				break
			}
			l.Stats <- &liveStats
		}
	}()

	<-l.stop
}

func (l *Listener) Stop() {
	close(l.stop)
}
