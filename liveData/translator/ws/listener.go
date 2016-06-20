package ws

import (
	"encoding/json"
	"fmt"
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
	stop   chan struct{}
}

func NewListener(address string, path string, opts map[string]string) *Listener {
	return &Listener{
		address: address,
		path:    path,
		opts:    opts,
		Errors:  make(chan error),
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
				close(l.stop)
				return
			}
			var liveStats schema.LiveStats
			err = json.Unmarshal(message, &liveStats)
			fmt.Println(err)
			for _, game := range liveStats {
				for _, player := range game.PlayerStats {
					fmt.Printf("%+v\n", player)
				}
			}
		}
	}()

	<-l.stop
}

func (l *Listener) Stop() {
	close(l.stop)
}
