package logs

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
)

type logString struct {
	message string
}

type ctxLog struct {
	sync.Mutex
	data map[string]*logString
}

type key int

const logsKey key = 1

func logContext(ctx context.Context, path string) {
	logs, ok := ctx.Value(logsKey).(*ctxLog)
	if !ok {
		return
	}
	buf := bytes.NewBufferString(path)
	for _, value := range logs.data {
		buf.WriteString(fmt.Sprintf("\n\t%s", value.message))
	}
	fmt.Println(buf.String())
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ctx = context.WithValue(ctx, logsKey, &ctxLog{
			data: make(map[string]*logString),
		})
		defer logContext(ctx, r.URL.Path)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
