package counter

import (
	"time"
)

func Counter(threshold time.Duration, resetCh <-chan bool, incrementCh <-chan time.Duration, triggerCh chan bool) {
	count := time.Second * 0
	for {
		select {
			case <-resetCh:
				count = 0
			case inc := <-incrementCh:
				count += inc
			default:
				if count >= threshold {
					triggerCh <- true
					count = 0
				}
				time.Sleep(time.Second/10)
		}
	}
}
