package channels

import (
	"fmt"
	"testing"
	"time"
)

func worker(c <-chan int) {
	select {
	case <-c:
		fmt.Println("do something...")
	case <-time.After(time.Second * 2):
		fmt.Println("timeout")
		return
	}
}

// channel 配合select 实现超时机制
func TestChannelWithSelectTimeout(t *testing.T) {
	c := make(chan int)
	worker(c)
}
