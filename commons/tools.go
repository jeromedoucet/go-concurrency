package commons

import (
	"sync"
	"time"
)

//todo units tests
func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

// todo unit tests
func WaitAnswerWithTimeOut(c chan bool, timeout time.Duration) (res bool, isTimedOut bool) {
	select {
	case res = <-c:
		return
	case <-time.After(timeout):
		isTimedOut = true
		return
	}
}
