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
	Stats  chan map[string]*schema.Game
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
		Stats:   make(chan map[string]*schema.Game),
		stop:    make(chan struct{}),
	}
}

func (l *Listener) Start() {
	go l.Run()
}

func (l *Listener) Run() {
	for {
		select {
		case <-l.stop:
			return
		default:
			l.connectAndListen()
		}
	}
}

func (l *Listener) connectAndListen() {
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

	for {
		select {
		case <-l.stop:
			return
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				l.Errors <- errors.Trace(err)
				return
			}
			var liveStats map[string]*schema.Game
			err = json.Unmarshal(message, &liveStats)
			if err != nil {
				l.Errors <- errors.Trace(err)
				break
			}
			if liveStats != nil {
				l.Stats <- liveStats
			}
		}
	}
}

func (l *Listener) Stop() {
	close(l.stop)
}
