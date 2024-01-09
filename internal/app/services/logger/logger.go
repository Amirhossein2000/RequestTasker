package logger

import "go.uber.org/zap"

type Logger struct {
	debug bool
	z     *zap.Logger
}

func NewLogger(debug bool) (*Logger, error) {
	z, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &Logger{
		debug: debug,
		z:     z,
	}, nil
}

func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Errorw(msg, keysAndValues...)
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Infow(msg, keysAndValues...)
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	if l.debug {
		l.z.Sugar().Debugw(msg, keysAndValues...)
	}
}
