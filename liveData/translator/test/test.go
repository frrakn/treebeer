package main

import "net"
import "fmt"
import "bufio"

func main() {

	for {
		// connect to this socket
		conn, err := net.Dial("tcp", "localhost:8082")
		if err != nil {
			fmt.Println(err)
			continue
		}
		for {
			message, _ := bufio.NewReader(conn).ReadString('}')
			fmt.Println("Message from server: " + message)
		}
	}
}
