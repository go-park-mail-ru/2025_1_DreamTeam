package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	MailRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mail_requests_total",
			Help: "Количество запросов на микросервис по отправки писем",
		},
		[]string{"method", "status"},
	)

	MailRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "mail_request_duration_seconds",
			Help:    "Длительность отправки письма",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)
)

func Init() {
	prometheus.MustRegister(MailRequestsTotal, MailRequestDuration)
}
