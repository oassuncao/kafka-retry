package message

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	Attempt           = "retry_attempts"
	Origin            = "retry_origin"
	DlqTopic          = "retry_dlq_topic"
	Topic             = "kafka_topic"
	Message           = "kafka_messageKey"
	ReceivedMessage   = "kafka_receivedMessageKey"
	ReceivedTimestamp = "kafka_receivedTimestamp"
)

func GetHeader(headers *[]kafka.Header, name string) *kafka.Header {
	for i := 0; i < len(*headers); i++ {
		if (*headers)[i].Key == name {
			return &(*headers)[i]
		}
	}
	return nil
}

func AddOrChange(msg *kafka.Message, name string, value []byte) bool {
	if header := GetHeader(&msg.Headers, name); header != nil {
		header.Value = value
		return false
	}

	header := kafka.Header{
		Key:   name,
		Value: value,
	}

	msg.Headers = append(msg.Headers, header)
	return true
}
