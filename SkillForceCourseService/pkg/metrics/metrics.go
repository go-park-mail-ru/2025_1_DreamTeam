package metrics

import (
	"net/http"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Init initializes Prometheus metrics server and registers default metrics.
func Init(port string) {
	// Регистрирует стандартные метрики gRPC
	grpc_prometheus.EnableHandlingTimeHistogram() // (по желанию)

	// Запуск HTTP-сервера с /metrics
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(port, nil); err != nil {
			panic("failed to start metrics server: " + err.Error())
		}
	}()
}
