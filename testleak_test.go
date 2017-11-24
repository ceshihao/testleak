package testleak

import (
	"testing"
	"time"
)

func TestLeakNoGoroutine(t *testing.T) {
	defer TestLeak(t)()
	func() {
		time.Sleep(5000 * time.Millisecond)
	}()
}

func TestLeakShortGoroutine(t *testing.T) {
	defer TestLeak(t)()
	go func() {
		time.Sleep(5 * time.Millisecond)
	}()
}

func TestLeakLongGoroutine(t *testing.T) {
	// demo Goroutine leak failure case by comment the following line
	t.Skip("Skip this case because Goroutine leak failure")
	defer TestLeak(t)()
	go func() {
		time.Sleep(5000 * time.Millisecond)
	}()
}

func TestLeakAppendWhiteList(t *testing.T) {
	AppendTestLeakWhiteList("time.Sleep(")
	defer TestLeak(t)()
	go func() {
		time.Sleep(5000 * time.Millisecond)
	}()
}
