package channels

import (
	"log"
	"sync"
	"testing"
	"time"
)

// 将带缓冲channel用作计数信号量
// 利用了带缓冲channel允许异步写入，并且缓冲区满的时候，会阻塞的特性

// 因为带缓冲channel的缓冲区大小就是信号量的值。
// 向带缓冲channel发送数据时，就相当于获取一个信号量
// 从带缓冲channel接收数据时，就相当于释放一个信号量
var active = make(chan struct{}, 3)
var jobs = make(chan int, 10)

func TestBufferedChannelAsCountingSemaphore(t *testing.T) {
	go func() {
		for i := 0; i < 8; i++ {
			jobs <- i + 1
		}
		close(jobs)
	}()

	var wg sync.WaitGroup
	for j := range jobs {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()

			active <- struct{}{}
			log.Printf("handle job: %d\n", j)
			time.Sleep(2 * time.Second)
			<-active
		}(j)
	}

	wg.Wait()

	// 2024/04/23 22:52:35 handle job: 8
	//2024/04/23 22:52:35 handle job: 2
	//2024/04/23 22:52:35 handle job: 1
	//2024/04/23 22:52:37 handle job: 3
	//2024/04/23 22:52:37 handle job: 7
	//2024/04/23 22:52:37 handle job: 6
	//2024/04/23 22:52:39 handle job: 4
	//2024/04/23 22:52:39 handle job: 5

	// 从运行结果时间戳可以看到，只有三个协程同时执行，其他协程被阻塞在active <- struct{}{}处
	// 缓冲区为3，所以当缓冲区满时，会阻塞发送，同一时间处于活跃状态(处理job)的协程数量最多为3
}
