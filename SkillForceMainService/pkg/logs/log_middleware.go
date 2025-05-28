package logs

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type LogString struct {
	Message string
}

type CtxLog struct {
	sync.Mutex
	Data []*LogString
}

type key int

const LogsKey key = 1

func logContext(ctx context.Context, path string, start time.Time) {
	logs, ok := ctx.Value(LogsKey).(*CtxLog)
	if !ok {
		return
	}

	duration := time.Since(start)

	buf := bytes.NewBufferString(path)
	buf.WriteByte('\n')
	fmt.Fprintf(buf, "Request duration: %v\n", duration)

	for _, log := range logs.Data {
		fmt.Fprintf(buf, "\t%s", log.Message)
	}

	fmt.Println(buf.String())
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		start := time.Now()

		ctx = context.WithValue(ctx, LogsKey, &CtxLog{
			Data: make([]*LogString, 0),
		})
		defer logContext(ctx, r.URL.Path, start)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
