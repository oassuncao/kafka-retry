package kafka

import (
	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/oassuncao/kafka-retry/pkg/config"
	"go.uber.org/zap"
)

type Producer struct {
	producer   *confluent.Producer
	ErrorEvent chan error
}

func (producer *Producer) Close() {
	close(producer.ErrorEvent)
	producer.producer.Close()
}

func (producer *Producer) SendMessage(message *confluent.Message) error {
	return producer.SendMessageDelivery(message, nil)
}

func (producer *Producer) SendMessageDelivery(message *confluent.Message, eventChan chan confluent.Event) error {
	return producer.producer.Produce(message, eventChan)
}

func NewProducer(kafkaConf config.Kafka) (*Producer, error) {
	configMap := config.LoadProducer(kafkaConf)
	producer, err := confluent.NewProducer(&configMap)
	if err != nil {
		return nil, err
	}

	producerResult := Producer{producer: producer}
	producerResult.ErrorEvent = make(chan error)
	// Delivery report handler for produced messages
	go func() {
		for e := range producer.Events() {
			switch event := e.(type) {
			case *confluent.Message:
				if event.TopicPartition.Error != nil {
					zap.S().Error("Delivery failed: ", event.TopicPartition)
					producerResult.ErrorEvent <- event.TopicPartition.Error
				}
			}
		}
	}()
	return &producerResult, nil
}
