// Go 1.2
// go run helloworld_go.go

package main

import (
    . "fmt"
    "runtime"
    "time"
)

var i = 0

func someGoroutine1() {
    for j := 0; j < 1000000; j++ {
	i++
    }
}

func someGoroutine2() {
    for j := 0; j < 1000000; j++ {
	i--
    }
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    go someGoroutine1() 
    go someGoroutine2() 

    time.Sleep(100*time.Millisecond)
    Println(i)
}
