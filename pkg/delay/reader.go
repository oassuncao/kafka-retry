package delay

import (
	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/oassuncao/kafka-retry/pkg/config"
	"github.com/oassuncao/kafka-retry/pkg/kafka"
	"github.com/oassuncao/kafka-retry/pkg/message"
	"go.uber.org/zap"
	"time"
)

type Reader struct {
	duration time.Duration
	producer *kafka.Producer
	metric   *prometheus.Summary
	closed   bool
}

func NewReader(topic string, duration time.Duration) *Reader {
	reader := Reader{duration: duration}
	reader.metric = kafka.CreateMetric(topic)

	producer, err := kafka.NewProducer(config.GetContext().Config.Kafka)
	if err != nil {
		zap.S().Error("Error on creating retry producer: ", err)
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

	zap.S().Infow("[Delay] - Received message", "consumer", consumer.Name, "topic", *msg.TopicPartition.Topic, "key", string(msg.Key))

	if duration := message.GetSleepTime(message.GetReceivedTime(msg), reader.duration); duration > 0 {
		time.Sleep(duration)
	}

	reader.observeMetric(msg)

	newMessage, err := message.NewOriginMessage(msg)
	if err != nil {
		zap.S().Errorf("Error on creating new message %s %v", consumer.Topics, err)
		return
	}

	zap.S().Infow("[Delay] - Sending message", "consumer", consumer.Name, "topic", *msg.TopicPartition.Topic, "key", string(msg.Key), "sendTopic", *newMessage.TopicPartition.Topic)
	err = reader.producer.SendMessage(newMessage)
	if err != nil {
		zap.S().Errorf("Error on sending message with key %s to topic %s %v", string(newMessage.Key), consumer.Topics, err)
	}
}

func (reader *Reader) observeMetric(msg *confluent.Message) {
	if reader.metric == nil {
		return
	}
	elapsedTime := message.GetElapsedTime(message.GetReceivedTime(msg), reader.duration)
	(*reader.metric).Observe(elapsedTime.Seconds())
}

func (reader *Reader) Close() {
	if reader.closed {
		return
	}

	reader.closed = true
	reader.producer.Close()
}
