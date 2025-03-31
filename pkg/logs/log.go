package logs

import (
	"bytes"
	"context"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logrusLogger *logrus.Logger
	buf          bytes.Buffer
}

func NewLogger() *Logger {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)

	return &Logger{
		logrusLogger: logger,
		buf:          buf,
	}
}

func (l *Logger) PrintLog(ctx context.Context, funcName string, message string) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)

	logger.WithField("function", funcName).Info(message)
	ctxLog, _ := ctx.Value(LogsKey).(*CtxLog)

	ctxLog.Data[funcName] = &LogString{
		Message: l.buf.String(),
	}
}
