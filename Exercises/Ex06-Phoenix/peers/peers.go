package peers

import (
	"time"
	"./udp"
)

func Peer(port int, sendCh <-chan string, receiveCh chan string) {
	go Transmitter(port, sendCh)
	go Receiver(port, receiveCh)
}

func Transmitter(port int, sendCh <-chan string) {
	c := udp.Init(port)
	for {
		select {
			case msg := <- sendCh:
				udp.Broadcast(c, msg)
			default: 
				time.Sleep(50*time.Millisecond)
		}
	}
}

func Receiver(port int, receiveCh chan string) {
	c := udp.Init(port)
	for {
		receiveCh <- udp.Receive(c)
		time.Sleep(50*time.Millisecond)
	}
}
