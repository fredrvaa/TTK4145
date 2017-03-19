package timer

import (
	"time"
)

func Timer(duration time.Duration, timeoutCh chan bool) {
	for {
		time.Sleep(duration)
		timeoutCh <- true
	}
}
