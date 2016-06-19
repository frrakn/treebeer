/*
===========
Connections
===========
Provides TCP / websocket management

*/

package connections

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var (
	newConnections      map[*websocket.Conn]struct{}
	existingConnections map[*websocket.Conn]struct{}
	upgrader            websocket.Upgrader
)

const (
	network = "tcp"
)

func init() {
	newConnections = map[*websocket.Conn]struct{}{}
	existingConnections = map[*websocket.Conn]struct{}{}
	upgrader = websocket.Upgrader{
		// eventually check the origin for this
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func OpenTcp(address string) chan map[string]int {
	address_tcp, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		panic(fmt.Sprintf("Connections: Unable to resolve TCP address %s\n", err))
	}

	connection, err := net.DialTCP(network, nil, address_tcp)
	if err != nil {
		panic(fmt.Sprintf("Connections: unable to create TCP listener %s\n", err))
	}

	connOutput := make(chan map[string]int)

	go func() {
		defer connection.Close()
		decoder := json.NewDecoder(connection)
		for {
			var message map[string]int
			err := decoder.Decode(&message)
			if err != nil {
				fmt.Printf("Connections: error parsing into json: %s\n", err)
				return
			}
			connOutput <- message
		}
	}()

	return connOutput
}

func HandleNewWsConnection(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Connections: unable to upgrade websocket connection from %s: %s\n", r.Referer(), err)
	}

	go func() {
		defer closeConn(conn)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				fmt.Sprintf("Connections: Error [%s] at %s, closing socket\n", err, conn.RemoteAddr())
				break
			}
		}
	}()

	newConnections[conn] = struct{}{}
}

func PromoteConns(message map[string]int) {
	for conn, _ := range newConnections {
		promoteConn(conn, message)
	}
}

func promoteConn(c *websocket.Conn, message map[string]int) {
	err := sendMessage(c, message)
	if err != nil {
		fmt.Printf("Connections: Failed to promote connection %s: %s", c, err)
	}
	delete(newConnections, c)
	existingConnections[c] = struct{}{}
}

func sendMessage(c *websocket.Conn, m map[string]int) error {
	json, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = c.WriteMessage(websocket.TextMessage, json)
	if err != nil {
		return err
	}

	return nil
}

func Close() {
	closePool(newConnections)
	closePool(existingConnections)
}

func closePool(pool map[*websocket.Conn]struct{}) {
	for conn, _ := range pool {
		closeConn(conn)
	}
}

func closeConn(c *websocket.Conn) {
	c.Close()
	delete(newConnections, c)
	delete(existingConnections, c)
}

func Broadcast(message map[string]int) error {
	var err error
	for conn, _ := range existingConnections {
		failed := sendMessage(conn, message)
		if failed != nil {
			err = failed
		}
	}
	return err
}
