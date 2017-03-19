package main

import (
	"fmt"
	"time"
	"./peers"
)

func main() {

	// Make threads
	sendCh := make(chan string)
	go peers.Transmitter(20009, sendCh)

	//conn := udp.Init(20009)
	// Main loop
	for {
		sendCh <- "test"
		fmt.Println("sent")
		time.Sleep(time.Second/10)
	}	
}
