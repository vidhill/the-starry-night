package stubrepository

type StubLogger struct{}

func (l *StubLogger) Debug(v ...interface{}) {}
func (l *StubLogger) Info(v ...interface{})  {}
func (l *StubLogger) Warn(v ...interface{})  {}
func (l *StubLogger) Error(v ...interface{}) {}

func NewStubLogger() StubLogger {
	return StubLogger{}
}
