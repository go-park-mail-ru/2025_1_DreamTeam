package logs

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
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
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf("Request duration: %v\n", duration))
	for _, log := range logs.Data {
		buf.WriteString(fmt.Sprintf("\t%s", log.Message))
	}
	fmt.Println(buf.String())
}

func GRPCLoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()

		ctx = context.WithValue(ctx, LogsKey, &CtxLog{
			Data: make([]*LogString, 0),
		})

		resp, err = handler(ctx, req)

		logContext(ctx, info.FullMethod, start)
		return
	}
}
