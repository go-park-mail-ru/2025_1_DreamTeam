package logs

import (
	"bytes"
	"context"
	"fmt"

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
	l.logrusLogger.WithField("function", funcName).Info(message)
	ctxLog, _ := ctx.Value(logsKey).(*ctxLog)

	ctxLog.Lock()
	defer ctxLog.Unlock()
	defer l.buf.Reset()

	if _, exist := ctxLog.data[funcName]; exist {
		ctxLog.data[funcName].message += fmt.Sprintf("\n\t%s", l.buf.String())
		return
	}

	ctxLog.data[funcName] = &logString{
		message: l.buf.String(),
	}

}
