package timer

import (
	"time"
)

func Timer(duration time.Duration, resetCh chan bool, timeoutCh chan bool) {
	var d time.Duration = 0
	for {
		select {
		case reset := <-resetCh:
			if reset {
				d = duration
			} else {
				d = 0
			}
		case <-time.After(100*time.Millisecond):
			if (d > 0) {
				if (d > 100*time.Millisecond) {
					d -= 100*time.Millisecond
				} else {
					d = 0
					timeoutCh <- true
				}
			}
		}
	}
}
