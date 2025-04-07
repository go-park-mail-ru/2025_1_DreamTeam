package logs

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
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

func logContext(ctx context.Context, path string) {
	logs, ok := ctx.Value(LogsKey).(*CtxLog)
	if !ok {
		return
	}

	buf := bytes.NewBufferString(path)
	buf.WriteString("\n")
	for _, log := range logs.Data {
		buf.WriteString(fmt.Sprintf("\t%s", log.Message))
	}
	fmt.Println(buf.String())
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ctx = context.WithValue(ctx, LogsKey, &CtxLog{
			Data: make([]*LogString, 0),
		})
		defer logContext(ctx, r.URL.Path)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
