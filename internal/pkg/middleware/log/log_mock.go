package log

type LoggerMock interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	DPanicw(msg string, keysAndValues ...interface{})
	Panic(args ...interface{})
	Panicf(template string, args ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Sync() error
}

func NewLoggerMock() LoggerMock {
	return &loggerMock{}
}

type loggerMock struct{}

func (l *loggerMock) Debug(args ...interface{})                        {}
func (l *loggerMock) Debugf(template string, args ...interface{})      {}
func (l *loggerMock) Debugw(msg string, keysAndValues ...interface{})  {}
func (l *loggerMock) Info(args ...interface{})                         {}
func (l *loggerMock) Infof(template string, args ...interface{})       {}
func (l *loggerMock) Infow(msg string, keysAndValues ...interface{})   {}
func (l *loggerMock) Warn(args ...interface{})                         {}
func (l *loggerMock) Warnf(template string, args ...interface{})       {}
func (l *loggerMock) Warnw(msg string, keysAndValues ...interface{})   {}
func (l *loggerMock) Error(args ...interface{})                        {}
func (l *loggerMock) Errorf(template string, args ...interface{})      {}
func (l *loggerMock) Errorw(msg string, keysAndValues ...interface{})  {}
func (l *loggerMock) DPanic(args ...interface{})                       {}
func (l *loggerMock) DPanicf(template string, args ...interface{})     {}
func (l *loggerMock) DPanicw(msg string, keysAndValues ...interface{}) {}
func (l *loggerMock) Panic(args ...interface{})                        {}
func (l *loggerMock) Panicf(template string, args ...interface{})      {}
func (l *loggerMock) Panicw(msg string, keysAndValues ...interface{})  {}
func (l *loggerMock) Fatal(args ...interface{})                        {}
func (l *loggerMock) Fatalf(template string, args ...interface{})      {}
func (l *loggerMock) Fatalw(msg string, keysAndValues ...interface{})  {}
func (l *loggerMock) Sync() error                                      { return nil }
