package main

import (
	"fmt"
	"time"
	"./peers"
)

func main() {

	// Make threads
	receiveCh := make(chan string)
	go peers.Receiver(20009, receiveCh)

	// Main loop
	for {
		select {
			case msg := <-receiveCh:
				fmt.Println(msg)
			default:
				time.Sleep(time.Second/10)

		}
	}
}
