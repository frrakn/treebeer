// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:9721", "http service address")

func main() {
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	header := http.Header{}
	header.Add("Host", "livestats.proxy.lolesports.com")
	header.Add("Pragma", "no-cache")
	header.Add("Cache-Control", "no-cache")
	header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.116 Safari/537.36")
	header.Add("Accept-Encoding", "gzip, deflate, sdch")
	header.Add("Accept-Language", "en-US,en;q=0.8")
	header.Add("Sec-WebSocket-Extensions", "permessage-deflate; client_max_window_bits")

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	/*db, err := sql.Open("mysql",
		"root:@tcp(127.0.0.1:3306)/treebeer")
	if err != nil {
		log.Fatal(err)
	}*/

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)

			/*			stmt, err := db.Prepare("INSERT INTO event(created_at, raw_message) VALUES(?, ?)")
						if err != nil {
							log.Fatal(err)
						}
						res, err := stmt.Exec(time.Now(), message)
						if err != nil {
							log.Fatal(err)
						}
						lastId, err := res.LastInsertId()
						if err != nil {
							log.Fatal(err)
						}
						rowCnt, err := res.RowsAffected()
						if err != nil {
							log.Fatal(err)
						}

						log.Printf("lastId: %s", lastId)
						log.Printf("rowCnt: %s", rowCnt)*/
		}
	}()

	<-interrupt
}
