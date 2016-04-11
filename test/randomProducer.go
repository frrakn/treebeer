/*
===============
RANDOM PRODUCER
===============
Accepts a single TCP connection, keeps a map with keys of 'a' to 'z' to integers initialized to 0, and then every second, randomly generates a delta, which is broadcasted and added to the in-memory map, and prints the in-memory map

Usage: Testing LiveDataAggregator

*/

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
)

var (
	address = flag.String("address", "localhost:9721", "http address")
)

const (
	network     = "tcp"
	maxNumKeys  = 3
	maxKeyValue = 10
	maxInc      = 10
)

type Producer struct {
	totalRandos   map[string]int
	randosToTrack chan map[string]int
	randosToSend  chan map[string]int
	writer        *bufio.Writer
}

func main() {
	flag.Parse()
	fmt.Printf("Main: Starting Random Producer at %s...\n", *address)

	address_tcp, _ := net.ResolveTCPAddr(network, *address)
	listener, _ := net.ListenTCP(network, address_tcp)

	connection, _ := listener.AcceptTCP()
	fmt.Println("Main: Connection started!")

	defer connection.Close()

	producer := new(Producer)
	producer.totalRandos = make(map[string]int)
	producer.randosToSend = make(chan map[string]int)
	producer.randosToTrack = make(chan map[string]int)
	producer.writer = bufio.NewWriter(connection)

	fmt.Println("Main: Starting to rando generation...")
	go producer.generateRandos(maxNumKeys, maxKeyValue, maxInc)

	fmt.Println("Main: Starting to track randos...")
	go producer.trackRandos()

	fmt.Println("Main: Starting to send off randos...")
	go producer.sendRandos()

	forever := make(chan bool)
	<-forever
}

func (p *Producer) generateRandos(maxNumKeys int, maxKeyValue int, maxInc int) {
	for {
		time.Sleep(1000 * time.Millisecond)
		rando := make(map[string]int)
		numGen := rand.Intn(maxNumKeys)

		for i := 0; i <= numGen; i++ {
			key, inc := generateKeyValue(maxKeyValue, maxInc)
			rando[strconv.Itoa(key)] += inc
		}

		fmt.Printf("Generate: Generated rando %+v\n", rando)

		p.randosToSend <- rando
		p.randosToTrack <- rando
	}
}

func (p *Producer) trackRandos() {
	for {
		rando := <-p.randosToTrack

		for key, value := range rando {
			p.totalRandos[key] += value
		}
		fmt.Printf("Track: Total randos updated to %+v with rando %+v\n", p.totalRandos, rando)
	}
}

func (p *Producer) sendRandos() {
	for {
		rando := <-p.randosToSend

		stringRando := make(map[string]string)
		for key, value := range rando {
			stringRando[key] = strconv.Itoa(value)
		}

		randoBytes, _ := json.Marshal(stringRando)
		p.writer.Write(randoBytes)
		p.writer.Flush()
	}
}

func generateKeyValue(maxKeyValue int, maxInc int) (int, int) {
	return rand.Intn(maxKeyValue), rand.Intn(maxInc) + 1
}
