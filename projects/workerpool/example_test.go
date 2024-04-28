package workerpool

import (
	"testing"
	"time"
)

func TestWorkerpool(t *testing.T) {
	pool := New(5)

	for i := 0; i < 10; i++ {
		err := pool.Schedule(func() {
			time.Sleep(time.Second * 3)
		})
		if err != nil {
			println("task: ", i, "err:", err)
		}
	}

	pool.Free()
}
