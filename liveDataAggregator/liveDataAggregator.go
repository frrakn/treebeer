/*
====================
LIVE DATA AGGREGATOR
====================
Provides data aggregation service for incoming map[string]int

*/

package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/frrakn/treebeer/liveDataAggregator/aggregator"
	"github.com/frrakn/treebeer/liveDataAggregator/connections"

	"github.com/julienschmidt/httprouter"
)

var (
	translatorAddress      = flag.String("translator_address", "localhost:9721", "translator tcp socket address")
	promotionInterval      = flag.Int("promote", 2, "interval in seconds at which connections get promoted")
	listeningPort          = flag.Int("listen_port", 9720, "address at which server listens for new websocket connections")
	broadcastInterval      = flag.Int("broadcast", 1, "interval in seconds between broadcasts")
	prelimAggregator       = aggregator.NewAggregator()
	fullAggregator         = aggregator.NewAggregator()
	promotionBroadcastSema = make(chan struct{}, 1)
)

func main() {
	flag.Parse()
	fmt.Printf("LiveDataAggregator: starting...\n")

	translator := connections.OpenTcp(*translatorAddress)
	go handleTranslator(translator)

	router := httprouter.New()
	router.GET("/", connections.HandleNewWsConnection)
	go http.ListenAndServe(":"+strconv.Itoa(*listeningPort), router)

	//todo(fchen): consider getting better synchronization, wait on new connections for broadcast
	promotionBroadcastSema <- struct{}{}
	go promote()
	go broadcast()

	//todo(fchen): kill by os signal ad graceful shutdown
	forever := make(chan struct{})
	<-forever
}

func handleTranslator(translator chan map[string]int) {
	for {
		newMessage := <-translator
		prelimAggregator.Add(newMessage)
	}
}

func promote() {
	for {
		<-promotionBroadcastSema
		message := fullAggregator.GetSum()
		connections.PromoteConns(message)
		fmt.Printf("Promoting %+v\n", message)
		promotionBroadcastSema <- struct{}{}
		time.Sleep(time.Duration(*promotionInterval) * time.Second)
	}
}

func broadcast() {
	for {
		<-promotionBroadcastSema
		message := prelimAggregator.GetSum()
		prelimAggregator.Zero()
		connections.Broadcast(message)
		fullAggregator.Add(message)
		fmt.Printf("Broadcasting %+v\n", message)
		promotionBroadcastSema <- struct{}{}
		time.Sleep(time.Duration(*broadcastInterval) * time.Second)
	}
}
