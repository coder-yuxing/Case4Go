package instrumenter

// Instrumenter 自动注入接口
type Instrumenter interface {
	Instrument(string) ([]byte, error)
}
