package service

import (
	"github.com/oassuncao/kafka-retry/pkg/config"
	"github.com/oassuncao/kafka-retry/pkg/delay"
	"github.com/oassuncao/kafka-retry/pkg/kafka"
	"github.com/oassuncao/kafka-retry/pkg/retry"
	"go.uber.org/zap"
	"sync"
)

type Manager struct {
	retries  []*kafka.Consumer
	attempts []*kafka.Consumer
}

func NewManager() *Manager {
	manager := Manager{
		retries:  []*kafka.Consumer{},
		attempts: []*kafka.Consumer{},
	}
	return &manager
}

func (manager *Manager) Start() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go startRetry(manager, &wg)

	wg.Add(1)
	go startAttempt(manager, &wg)

	wg.Wait()
}

func (manager *Manager) Stop() {
	wg := sync.WaitGroup{}
	zap.S().Info("Stopping services...")
	for _, item := range manager.retries {
		wg.Add(1)
		go func(item *kafka.Consumer) {
			item.Stop()
			wg.Done()
		}(item)
	}

	for _, item := range manager.attempts {
		wg.Add(1)
		go func(item *kafka.Consumer) {
			item.Stop()
			wg.Done()
		}(item)
	}
	wg.Wait()
	zap.S().Info("Stopped all services...")
}

func (manager *Manager) Close() {
	zap.S().Info("Closing services...")
	for _, item := range manager.retries {
		item.Close()
	}

	for _, item := range manager.attempts {
		item.Close()
	}
	zap.S().Info("Closed services...")
}

func (manager *Manager) GetRetriesSize() int {
	return len(manager.retries)
}

func (manager *Manager) GetAttemptsSize() int {
	return len(manager.attempts)
}

func startAttempt(manager *Manager, wg *sync.WaitGroup) {
	context := config.GetContext()
	consumerConfig := config.NewConsumer(context.Config.Retry.Group)
	for _, attempt := range context.Config.Retry.Attempts {
		reader := delay.NewReader(attempt.Topic, attempt.Delay)
		for i := 0; i < attempt.Concurrency; i++ {
			wg.Add(1)

			consumer, err := startConsumer(consumerConfig, attempt.Topic, reader)
			if err != nil {
				zap.S().Error("Error on creating attempt consumer ", err)
			} else {
				manager.attempts = append(manager.attempts, consumer)
			}

			wg.Done()
		}
	}
	wg.Done()
}

func startRetry(manager *Manager, wg *sync.WaitGroup) {
	context := config.GetContext()
	consumerConfig := config.NewConsumer(context.Config.Retry.Group)
	reader := retry.NewReader(context.Config.Retry.Topic)
	for i := 0; i < context.Config.Retry.Concurrency; i++ {
		wg.Add(1)

		consumer, err := startConsumer(consumerConfig, context.Config.Retry.Topic, reader)
		if err != nil {
			zap.S().Error("Error on creating retry consumer ", err)
		} else {
			manager.retries = append(manager.retries, consumer)
		}

		wg.Done()
	}
	wg.Done()
}

func startConsumer(consumerConfig config.Consumer, topic string, reader kafka.Reader) (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(config.GetContext().Config.Kafka, consumerConfig)
	if err != nil {
		return nil, err
	}

	err = consumer.Subscribe(topic)
	if err != nil {
		return nil, err
	}

	go consumer.ReadMessage(reader)
	return consumer, nil
}
