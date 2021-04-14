package kafka

import (
	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/oassuncao/kafka-retry/pkg/config"
	"go.uber.org/zap"
)

type Reader interface {
	Read(Consumer, *confluent.Message, error)
	Close()
}

func CreateMetric(topic string) *prometheus.Summary {
	context := config.GetContext()
	summary := prometheus.NewSummary(prometheus.SummaryOpts{
		Name: context.Config.Retry.Metric.Name,
		ConstLabels: prometheus.Labels{
			"topic": topic,
		},
		Help: context.Config.Retry.Metric.Description,
	})

	err := prometheus.Register(summary)
	if err != nil {
		zap.S().Errorf("Error on register metric ", err)
	}
	return &summary
}
