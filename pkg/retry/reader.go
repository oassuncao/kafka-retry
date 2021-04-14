package retry

import (
	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/oassuncao/kafka-retry/pkg/config"
	"github.com/oassuncao/kafka-retry/pkg/kafka"
	"github.com/oassuncao/kafka-retry/pkg/message"
	"go.uber.org/zap"
)

type Reader struct {
	producer *kafka.Producer
	closed   bool
	metric   *prometheus.Summary
}

func NewReader(topic string) *Reader {
	reader := Reader{}
	reader.metric = kafka.CreateMetric(topic)

	producer, err := kafka.NewProducer(config.GetContext().Config.Kafka)
	if err != nil {
		zap.S().Errorf("Error on creating retry producer ", err)
	} else {
		reader.producer = producer
	}
	return &reader
}

func (reader *Reader) Read(consumer kafka.Consumer, msg *confluent.Message, err error) {
	if err != nil {
		zap.S().Errorf("Error on reading message on topic %s %v", consumer.Topics, err)
		return
	}

	zap.S().Infow("[Retry] - Received message", "consumer", consumer.Name, "topic", *msg.TopicPartition.Topic, "key", string(msg.Key))

	reader.observeMetric(msg)
	newMessage, err := message.NewRetryMessage(msg, config.GetContext().Config.Retry)
	if err != nil {
		zap.S().Errorf("Error on creating new message %s %v", consumer.Topics, err)
		return
	}

	zap.S().Infow("[Retry] - Sending message", "consumer", consumer.Name, "topic", *msg.TopicPartition.Topic,
		"key", string(newMessage.Key), "sendTopic", *newMessage.TopicPartition.Topic)

	err = reader.producer.SendMessage(newMessage)
	if err != nil {
		zap.S().Errorf("Error on sending message with key %s to topic %s %v", string(newMessage.Key), consumer.Topics, err)
	}
}

func (reader *Reader) Close() {
	if reader.closed {
		return
	}

	reader.closed = true
	reader.producer.Close()
}

func (reader *Reader) observeMetric(msg *confluent.Message) {
	if reader.metric == nil {
		return
	}

	duration := message.GetElapsedTime(message.GetReceivedTime(msg), 0)
	(*reader.metric).Observe(duration.Seconds())
}
