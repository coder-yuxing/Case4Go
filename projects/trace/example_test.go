package trace_test

import (
	"yuxing.com/trace"
)

func A1() {
	defer trace.Trace()()
	B1()
}

func B1() {
	defer trace.Trace()()
	C1()
}

func C1() {
	defer trace.Trace()()
	D()
}

func D() {
	defer trace.Trace()()
}

func ExampleTrace() {
	A1()
	// Output:
	// g[00001]: ->yuxing.com/trace_test.A1
	// g[00001]:   ->yuxing.com/trace_test.B1
	// g[00001]:     ->yuxing.com/trace_test.C1
	// g[00001]:       ->yuxing.com/trace_test.D
	// g[00001]:         <-yuxing.com/trace_test.D
	// g[00001]:       <-yuxing.com/trace_test.C1
	// g[00001]:     <-yuxing.com/trace_test.B1
	// g[00001]:   <-yuxing.com/trace_test.A1
}
