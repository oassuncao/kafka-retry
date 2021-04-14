package message

import (
	"errors"
	"fmt"
	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/oassuncao/kafka-retry/pkg/config"
	"strconv"
	"time"
)

func NewRetryMessage(msg *confluent.Message, retry config.Retry) (*confluent.Message, error) {
	if GetHeader(&msg.Headers, Origin) == nil {
		return nil, errors.New(fmt.Sprintf("%s header was not provided or was empty", Origin))
	}

	attempts := getAttempts(msg) + 1
	topic := getTopicName(msg, retry, attempts)
	return newMessage(msg, topic, attempts)
}

func NewOriginMessage(msg *confluent.Message) (*confluent.Message, error) {
	header := GetHeader(&msg.Headers, Origin)
	if header == nil {
		return nil, errors.New(fmt.Sprintf("%s header was not provided or was empty", Origin))
	}
	return newMessage(msg, string(header.Value), getAttempts(msg))
}

func GetSleepTime(date time.Time, baseDuration time.Duration) time.Duration {
	return date.Add(baseDuration).Sub(time.Now().UTC())
}

func GetElapsedTime(date time.Time, baseDuration time.Duration) time.Duration {
	return time.Now().UTC().Sub(date.UTC().Add(baseDuration))
}

func GetReceivedTime(msg *confluent.Message) time.Time {
	date := time.Now().UTC()
	if msg.TimestampType != confluent.TimestampNotAvailable {
		date = msg.Timestamp.UTC()
	}
	return date
}

func newMessage(msg *confluent.Message, topic string, attempts int) (*confluent.Message, error) {
	newMessage := copyMessage(msg)
	newMessage.TopicPartition.Topic = &topic
	newMessage.Timestamp = time.Time{}

	AddOrChange(&newMessage, Attempt, []byte(strconv.Itoa(attempts)))
	return &newMessage, nil
}

func getTopicName(msg *confluent.Message, retry config.Retry, attempts int) string {
	dql := attempts > len(retry.Attempts)
	if dql {
		if header := GetHeader(&msg.Headers, DlqTopic); header != nil {
			return string(header.Value)
		}

		header := GetHeader(&msg.Headers, Origin)
		return fmt.Sprintf("%s.%s", string(header.Value), retry.Dlq)
	}
	return fmt.Sprintf("%s-%d", *msg.TopicPartition.Topic, attempts)
}

func copyMessage(msg *confluent.Message) confluent.Message {
	newMessage := *msg
	newMessage.TopicPartition.Partition = confluent.PartitionAny
	return newMessage
}

func getAttempts(msg *confluent.Message) int {
	if header := GetHeader(&msg.Headers, Attempt); header != nil {
		value, _ := strconv.Atoi(string(header.Value))
		return value
	}
	return 0
}
