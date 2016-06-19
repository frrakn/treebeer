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

var addr = flag.String("addr", "livestats.proxy.lolesports.com:80", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	header := http.Header{}
	header.Add("Host", "livestats.proxy.lolesports.com")
	header.Add("Pragma", "no-cache")
	header.Add("Cache-Control", "no-cache")
	header.Add("Origin", "http://fantasy.na.lolesports.com")
	header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.116 Safari/537.36")
	header.Add("Accept-Encoding", "gzip, deflate, sdch")
	header.Add("Accept-Language", "en-US,en;q=0.8")
	header.Add("Cookie", "__cfduid=d9eb189040e96fbcdbc1b7627632e37f41432410636; ping_session_id=988b83a1-c765-4b2d-955f-45a5a929857c; _ga=GA1.2.401925352.1432526363; PVPNET_LANG=en_US; PVPNET_TOKEN_NA=eyJkYXRlX3RpbWUiOjE0NTUzOTAyNDcsImdhc19hY2NvdW50X2lkIjoiMjIxNTkxNTkyIiwicHZwbmV0X2FjY291bnRfaWQiOiIyMjE1OTE1OTIiLCJzdW1tb25lcl9uYW1lIjoiRHI0Z09uWmJSZUFrRXIiLCJ2b3VjaGluZ19rZXlfaWQiOiI5MDM0NzUyYjJiNDU2MDQ0YWU4N2YyNTk4MmRhZDA3ZCIsInNpZ25hdHVyZSI6ImN5RUdZU2NUYXc2b0RUS21RL3BKS2x0RDFpSUtka05GREdtSEhwdFZBVjBWdHBnbHp0bkxTTVo1TGJHSVh4QThySTBreDRqYys0eW02VFZxeVFqV3VtZDFoY1pXTXBvRUJlTWE5dzFDWUcvMVpEM0lJQjdIYUFHOURyZEM0VVk1QVRhZm1LVlFUZUt2YjNvWW5qcm1Wd0ZUZ3Q2NHBEM1N4cVgzWE9WWGRHZz0ifQ%3D%3D; PVPNET_ACCT_NA=Dr4gOnZbReAkEr; PVPNET_ID_NA=221591592; PVPNET_REGION=na; ajs_user_id=%22na_221591592%22; ajs_group_id=null")
	header.Add("Sec-WebSocket-Extensions", "permessage-deflate; client_max_window_bits")

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/stats"}
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
