package main

import (
	"strconv"
	"fmt"
	"time"
	"./peers"
	"./timer"
	"./counter"
	"os/exec"
)

func spawnBackup() {

	cmd := exec.Command("gnome-terminal", "-x", "go", "run", "/home/vegst/elevator-lab-gr32/Ex06-Phoenix/phoenix.go")
	cmd.Start()
}

func main() {
	master := false
	i := 0

	// Make threads
	sendCh := make(chan string)
	receiveCh := make(chan string)
	go peers.Peer(20009, sendCh, receiveCh)

	timeoutCh := make(chan bool)
	go timer.Timer(time.Second, timeoutCh)

	resetCh := make(chan bool)
	incrementCh := make(chan time.Duration)
	triggerCh := make(chan bool)
	go counter.Counter(time.Second/2, resetCh, incrementCh, triggerCh)

	// Main loop
	fmt.Println("Now slave!")
	for {
		select {
			case <- timeoutCh:
				if master {
					i++
					fmt.Println(i)
				}
			case msg := <-receiveCh:
				if !master {
					resetCh <- true
					i,_ = strconv.Atoi(msg)
				}
			case <- triggerCh:
				master = true
				spawnBackup()
				fmt.Println("Now master!")
				
			default:
				if master {
					sendCh <- strconv.Itoa(i)
				}
				time.Sleep(time.Second/10)
				if !master {
					incrementCh <- time.Second/10
				}

		}
	}	
}
