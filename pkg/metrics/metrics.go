package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"hamsterbot/pkg/logger"
	"net/http"
)

var (
	IncomingMessages = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "telegram_bot_incoming_messages_total",
			Help: "Total number of incoming messages",
		},
	)
	MessageProcessingDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "telegram_bot_message_processing_duration_seconds",
			Help:    "Duration of message processing in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func Init() {
	prometheus.MustRegister(IncomingMessages)
	prometheus.MustRegister(MessageProcessingDuration)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		logger.Fatal("Ошибка при запуске HTTP-сервера", zap.Error(err))
	}
}
