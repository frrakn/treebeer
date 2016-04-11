/*
===========
Connections
===========
Provides TCP / websocket management

*/

package connections

import (
	"net"
)

const (
	network = "net"
)

func OpenTCP(address string) net.Conn {
	address_tcp, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		panic(fmt.Sprintf("Connections: Unable to resolve TCP address %s\n", err))
	}

	connection, err := net.DialTCP(network, nil, address_tcp)
	if err != nil {
		panic(fmt.Sprintf("Connections: unable to create TCP listener %s\n", err))
	}

	return connection
}
