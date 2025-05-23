package logs

import (
	"bytes"
	"context"

	"github.com/sirupsen/logrus"
)

func PrintLog(ctx context.Context, funcName string, message string) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)

	logger.WithField("function", funcName).Info(message)
	ctxLog, _ := ctx.Value(LogsKey).(*CtxLog)

	ctxLog.Lock()
	defer ctxLog.Unlock()

	ctxLog.Data = append(ctxLog.Data, &LogString{
		Message: buf.String(),
	})
}
