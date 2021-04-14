package kafka

import (
	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/oassuncao/kafka-retry/internal/util"
	"github.com/oassuncao/kafka-retry/pkg/config"
	"go.uber.org/zap"
)

type Consumer struct {
	consumer *confluent.Consumer
	Name     string
	Topics   []string
	stop     chan int
	stopped  chan int
}

func (consumer *Consumer) Close() {
	err := consumer.consumer.Close()
	if err != nil {
		zap.S().Errorf("Error on close consumer %s of topics %s %v", consumer.Name, consumer.Topics, err)
	}
}

func (consumer *Consumer) Stop() {
	close(consumer.stop)
	<-consumer.stopped
}

func (consumer *Consumer) Subscribe(topics ...string) error {
	consumer.Topics = topics
	return consumer.consumer.SubscribeTopics(topics, nil)
}

func (consumer *Consumer) ReadMessage(reader Reader) {
	for {
		select {
		case <-consumer.stop:
			reader.Close()
			err := consumer.consumer.Unsubscribe()
			if err != nil {
				zap.S().Errorf("Error on unsubscribe consumer %s of topics %s %v", consumer.Name, consumer.Topics, err)
			}
			close(consumer.stopped)
			return
		default:
			ev := consumer.consumer.Poll(1000)
			switch event := ev.(type) {
			case *confluent.Message:
				if event.TopicPartition.Error != nil {
					reader.Read(*consumer, event, event.TopicPartition.Error)
				} else {
					reader.Read(*consumer, event, nil)
				}
			case confluent.Error:
				reader.Read(*consumer, nil, event)
			default:
				if event == nil {
					continue
				}
			}
		}
	}
}

func NewConsumer(kafkaConf config.Kafka, consumerConf config.Consumer) (*Consumer, error) {
	configMap := config.LoadConsumer(kafkaConf, consumerConf)
	consumer, err := confluent.NewConsumer(&configMap)
	if err != nil {
		return nil, err
	}

	consumerResult := Consumer{Name: util.Random(15), consumer: consumer}
	consumerResult.stop = make(chan int)
	consumerResult.stopped = make(chan int)
	return &consumerResult, nil
}
