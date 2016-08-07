package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

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
		decoder := json.NewDecoder(conn)
		for {
			var message map[string]string
			err := decoder.Decode(&message)
			if err != nil {
				fmt.Println(err)
				continue
			}
			bytes, err := json.Marshal(message)

			fmt.Println(string(bytes))
			if _, err = f.Write(bytes); err != nil {
				fmt.Println(err)
			}
		}
	}
}
