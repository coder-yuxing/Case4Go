package channels

// channel 配合select 使用default避免阻塞
// 无论是无缓冲channel还是带缓冲channel都是适用的
// select 的default分支的语义就是在其他非default分支因通信未就绪，而无法被选择的时候，才会被执行。
// 这就给default分支赋予了一种“避免阻塞”的特性
func tyrRecv(c <-chan int) (int, bool) {
	select {
	case v := <-c:
		return v, true
	default:
		return 0, false
	}
}

// 在 Go 标准库中，这个方法的应用
// $GOROOT/src/time/sleep.go
//func sendTime(c interface{}, seq uintptr) {
//	// 无阻塞的向c发送当前时间
//	select {
//	case c.(chan Time) <- Now():
//	default:
//	}
//}
