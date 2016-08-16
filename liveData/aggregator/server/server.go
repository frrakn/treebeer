package server

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/frrakn/treebeer/liveData/aggregator/aggregator"
	pb "github.com/frrakn/treebeer/liveData/liveData"
	"github.com/gorilla/websocket"
	"github.com/juju/errors"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	config *Configuration

	router *httprouter.Router

	prelimAgg *aggregator.Aggregator
	finalAgg  *aggregator.Aggregator

	newConns         *connections
	existingConns    *connections
	promoteBroadcast sync.Mutex

	upgrader websocket.Upgrader

	stop   chan struct{}
	Errors chan error
}

type Configuration struct {
	Translator        string
	PromotionInterval int
	BroadcastInterval int
	Port              string
}

type connections struct {
	m map[*websocket.Conn]struct{}
	sync.Mutex
}

func NewServer(c *Configuration) *Server {
	s := &Server{
		config: c,

		router: httprouter.New(),

		prelimAgg: aggregator.NewAggregator(),
		finalAgg:  aggregator.NewAggregator(),

		newConns:      newConnections(),
		existingConns: newConnections(),

		upgrader: websocket.Upgrader{
			// TODO (frrakn): eventually check the origin for this
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},

		stop:   make(chan struct{}),
		Errors: make(chan error),
	}

	s.router.GET("/", s.handleWebsocket)

	return s
}

func (s *Server) Start() {
	go s.run()
}

func (s *Server) run() {
	go http.ListenAndServe(s.config.Port, s.router)
	go s.broadcastLoop()
	go s.promoteLoop()
	go s.listenTCP()
	<-s.stop
	s.newConns.close()
	s.existingConns.close()
}

func (s *Server) Stop() {
	close(s.stop)
}

func (s *Server) listenTCP() {
	conn, err := net.Dial("tcp", s.config.Translator)
	if err != nil {
		s.Errors <- errors.Trace(err)
		return
	}
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	for {
		var message *pb.Stat
		err := decoder.Decode(&message)
		if err != nil {
			s.Errors <- errors.Trace(err)
			continue
		}
		s.prelimAgg.AddProto(message)
	}
}

func (s *Server) handleWebsocket(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errors <- errors.Trace(err)
	}

	go func() {
		defer s.closeConn(conn)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				s.Errors <- errors.Trace(err)
				break
			}
		}
	}()

	s.newConns.add(conn)
}

func (s *Server) closeConn(c *websocket.Conn) {
	s.newConns.remove(c)
	s.existingConns.remove(c)
	c.Close()
}

func (s *Server) broadcast(agg *aggregator.Aggregator) {
	errs := s.existingConns.broadcast(agg)
	for _, err := range errs {
		s.Errors <- errors.Trace(err)
	}
}

func (s *Server) promoteLoop() {
	for {
		s.promoteBroadcast.Lock()
		s.promoteConns()
		s.promoteBroadcast.Unlock()
		time.Sleep(time.Duration(s.config.PromotionInterval) * time.Second)
	}
}

func (s *Server) promoteConns() {
	errs := s.newConns.broadcast(s.finalAgg)
	for _, err := range errs {
		s.Errors <- errors.Trace(err)
	}
	s.newConns.moveTo(s.existingConns)
}

func (s *Server) broadcastLoop() {
	for {
		s.promoteBroadcast.Lock()
		s.existingConns.broadcast(s.prelimAgg)
		s.finalAgg.AddProtos(s.prelimAgg.ToProtos())
		s.prelimAgg.Zero()
		s.promoteBroadcast.Unlock()
		time.Sleep(time.Duration(s.config.BroadcastInterval) * time.Second)
	}
}

func newConnections() *connections {
	return &connections{
		m: make(map[*websocket.Conn]struct{}),
	}
}

func (c *connections) remove(conn *websocket.Conn) {
	c.Lock()
	delete(c.m, conn)
	c.Unlock()
}

func (c *connections) add(conn *websocket.Conn) {
	c.Lock()
	c.m[conn] = struct{}{}
	c.Unlock()
}

func (c *connections) moveTo(c2 *connections) {
	c.Lock()
	c2.Lock()
	for conn, _ := range c.m {
		delete(c.m, conn)
		c2.m[conn] = struct{}{}
	}
	c2.Unlock()
	c.Unlock()
}

func (c *connections) close() {
	c.Lock()
	for conn, _ := range c.m {
		conn.Close()
	}
	c.Unlock()
}

func (c *connections) sendMessage(conn *websocket.Conn, agg *aggregator.Aggregator) error {
	protos := agg.ToProtos()

	for _, proto := range protos {
		msg, err := json.Marshal(proto)
		if err != nil {
			return errors.Trace(err)
		}
		err = conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

func (c *connections) broadcast(agg *aggregator.Aggregator) []error {
	var errs []error
	c.Lock()
	for conn, _ := range c.m {
		err := c.sendMessage(conn, agg)
		if err != nil {
			errs = append(errs, err)
		}
	}
	c.Unlock()
	return errs
}
