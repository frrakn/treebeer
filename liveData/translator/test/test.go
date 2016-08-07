package main

import (
	"fmt"
	"net"
	"os"
)

import "bufio"

func main() {

	for {
		f, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		// connect to this socket
		conn, err := net.Dial("tcp", "localhost:8082")
		if err != nil {
			fmt.Println(err)
			continue
		}
		for {
			message, _ := bufio.NewReader(conn).ReadString('}')
			fmt.Println(message)

			if _, err = f.WriteString(message); err != nil {
				fmt.Println(err)
			}
		}
	}
}
