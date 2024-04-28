package workerpool

import (
	"errors"
	"fmt"
	"sync"
)

var ErrWorkerPoolFreed = errors.New("workerpool freed") // workerpool已终止运行

const (
	defaultCapacity = 100  // 默认workerpool容量
	maxCapacity     = 1000 // workerpool最大容量
)

// Task 待执行任务
type Task func()

// Pool workerpool
type Pool struct {
	capacity int            // workerpool容量
	active   chan struct{}  // 活跃的worker数量
	tasks    chan Task      // 任务队列
	wg       sync.WaitGroup // 用于在pool销毁时等待所有worker的退出
	quit     chan struct{}  // 用于通知各个worker退出的信号channel
}

func (p *Pool) run() {
	idx := 0

	// run方法的无线循环，使用select监听Pool实例的两个channel: quit和active
	// 收到quit信号时，退出run方法
	// 收到active信号时，创建一个worker
	for {
		select {
		case <-p.quit:
			return
		case p.active <- struct{}{}:
			idx++
			p.newWorker(idx)
		}
	}
}

func (p *Pool) newWorker(i int) {
	p.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("worker[%03d]: recover panic[%s] and exit\n", i, err)
				<-p.active
			}
			p.wg.Done()
		}()

		fmt.Printf("worker[%03d]: start\n", i)
		for {
			select {
			case task := <-p.tasks:
				fmt.Printf("worker[%03d]: receive a task\n", i)
				task()
			case <-p.quit:
				fmt.Printf("worker[%03d]: exit\n", i)
				<-p.active
				return
			}
		}
	}()
}

func (p *Pool) Schedule(task Task) error {
	select {
	case <-p.quit:
		return ErrWorkerPoolFreed
	case p.tasks <- task:
		return nil
	}
}

func (p *Pool) Free() {
	close(p.quit)
	p.wg.Wait()
	fmt.Printf("workerpool freed\n")
}

// New 创建workerpool实例
func New(capacity int) *Pool {
	if capacity <= 0 {
		capacity = defaultCapacity
	}

	if capacity > maxCapacity {
		capacity = maxCapacity
	}

	p := &Pool{
		capacity: capacity,
		active:   make(chan struct{}, capacity),
		tasks:    make(chan Task),
		quit:     make(chan struct{}),
	}

	fmt.Printf("workerpool start\n")

	go p.run()

	return p
}
