package main

import (
	"fmt"

	"github.com/frrakn/treebeer/liveData/translator/ws"
)

func main() {
	config := &ws.Configuration{
		Address: "livestats.proxy.lolesports.com:80",
		Path:    "/stats",
		Opts: map[string]string{
			"Host":                     "livestats.proxy.lolesports.com",
			"Pragma":                   "no-cache",
			"Cache-Control":            "no-cache",
			"Origin":                   "http://fantasy.na.lolesports.com",
			"User-Agent":               "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.116 Safari/537.36",
			"Accept-Encoding":          "gzip, deflate, sdch",
			"Accept-Language":          "en-US,en;q=0.8",
			"Cookie":                   "__cfduid=d9eb189040e96fbcdbc1b7627632e37f41432410636; ping_session_id=988b83a1-c765-4b2d-955f-45a5a929857c; _ga=GA1.2.401925352.1432526363; PVPNET_LANG=en_US; PVPNET_TOKEN_NA=eyJkYXRlX3RpbWUiOjE0NTUzOTAyNDcsImdhc19hY2NvdW50X2lkIjoiMjIxNTkxNTkyIiwicHZwbmV0X2FjY291bnRfaWQiOiIyMjE1OTE1OTIiLCJzdW1tb25lcl9uYW1lIjoiRHI0Z09uWmJSZUFrRXIiLCJ2b3VjaGluZ19rZXlfaWQiOiI5MDM0NzUyYjJiNDU2MDQ0YWU4N2YyNTk4MmRhZDA3ZCIsInNpZ25hdHVyZSI6ImN5RUdZU2NUYXc2b0RUS21RL3BKS2x0RDFpSUtka05GREdtSEhwdFZBVjBWdHBnbHp0bkxTTVo1TGJHSVh4QThySTBreDRqYys0eW02VFZxeVFqV3VtZDFoY1pXTXBvRUJlTWE5dzFDWUcvMVpEM0lJQjdIYUFHOURyZEM0VVk1QVRhZm1LVlFUZUt2YjNvWW5qcm1Wd0ZUZ3Q2NHBEM1N4cVgzWE9WWGRHZz0ifQ%3D%3D; PVPNET_ACCT_NA=Dr4gOnZbReAkEr; PVPNET_ID_NA=221591592; PVPNET_REGION=na; ajs_user_id=%22na_221591592%22; ajs_group_id=null",
			"Sec-WebSocket-Extensions": "permessage-deflate; client_max_window_bits",
		},
	}

	listener := ws.NewListener(config)

	go handleListener(listener)

	listener.Run()
}

func handleListener(listener *ws.Listener) {
	for {
		select {
		case err := <-listener.Errors:
			fmt.Println(err)
		case stats := <-listener.Stats:
			fmt.Println(stats)
		}
	}
}
