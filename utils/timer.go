package utils

import (
	"time"
)

func StartTime(callback func(), interval int) (exitSignal chan struct{}) {
	exitSignal = make(chan struct{})
	go func() {
		for {
			select {
			case <-time.After(time.Duration(interval) * time.Millisecond):
				callback()
			case <-exitSignal:
				return
			}
		}
	}()
	return exitSignal
}
