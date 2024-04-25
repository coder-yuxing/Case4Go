package trace

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

var mu sync.Mutex
var m = make(map[uint64]int)

func Trace() func() {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("not found caller")
	}
	name := runtime.FuncForPC(pc).Name()

	gid := curGoroutineID()

	mu.Lock()
	indents := m[gid]    // 获取当前gid的缩进层次
	m[gid] = indents + 1 // 缩进层次+1后存入map
	mu.Unlock()

	printTrace(gid, name, "->", indents)
	return func() {
		mu.Lock()
		indents := m[gid]    // 获取当前gid对应的缩进层次
		m[gid] = indents - 1 // 缩进层次-1后存入map
		mu.Unlock()
		printTrace(gid, name, "<-", indents)
	}
}

func printTrace(gid uint64, name string, arrow string, indent int) {
	indents := ""
	for i := 0; i < indent; i++ {
		indents += "  "
	}
	fmt.Printf("g[%05d]: %s%s%s\n", gid, indents, arrow, name)
}

var goroutineSpace = []byte("goroutine ")

// curGoroutineID 获取当前协程ID
// 改造自Go标准库：$GOROOT/src/net/http/h2_bundle.go
func curGoroutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	// Parse the 4707 out of "goroutine 4707 ["
	b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		panic(fmt.Sprintf("No space found in %q", b))
	}
	b = b[:i]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse goroutine ID out of %q: %v", b, err))
	}
	return n
}
