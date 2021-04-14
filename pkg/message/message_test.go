package message

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/oassuncao/kafka-retry/pkg/config"
	"reflect"
	"testing"
	"time"
)

func TestGetElapsedTime(t *testing.T) {
	type args struct {
		date         time.Time
		baseDuration time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "normal",
			args: args{
				date:         time.Now(),
				baseDuration: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetElapsedTime(tt.args.date, tt.args.baseDuration); int(got.Seconds()) != int(tt.want.Seconds()) {
				t.Errorf("GetElapsedTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetReceivedTime(t *testing.T) {
	now := time.Now().UTC()
	type args struct {
		msg *kafka.Message
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "default",
			args: args{
				msg: &(kafka.Message{
					Timestamp:     now,
					TimestampType: kafka.TimestampCreateTime,
				}),
			},
			want: now,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetReceivedTime(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetReceivedTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSleepTime(t *testing.T) {
	type args struct {
		date         time.Time
		baseDuration time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "normal",
			args: args{
				date:         time.Now().Add(-(12 * time.Second)),
				baseDuration: 12 * time.Second,
			},
			want: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSleepTime(tt.args.date, tt.args.baseDuration); int(got.Seconds()) != int(tt.want.Seconds()) {
				t.Errorf("GetSleepTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOriginMessage(t *testing.T) {
	var topicName = "retry"
	type args struct {
		msg *kafka.Message
	}
	tests := []struct {
		name    string
		args    args
		want    *kafka.Message
		wantErr bool
	}{
		{
			name: "withoutOrigin",
			args: args{
				msg: &kafka.Message{},
			},
			wantErr: true,
		},
		{
			name: "withOrigin",
			args: args{
				msg: &kafka.Message{
					Headers: []kafka.Header{
						{
							Key:   Origin,
							Value: []byte(topicName),
						},
					},
				},
			},
			want: &kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &topicName,
					Partition: -1,
				},
				Timestamp:     time.Time{},
				TimestampType: 0,
				Headers: []kafka.Header{
					{
						Key:   Origin,
						Value: []byte(topicName),
					},
					{
						Key:   Attempt,
						Value: []byte("0"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOriginMessage(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOriginMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOriginMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRetryMessage(t *testing.T) {
	var topicName = "retry"
	var topicNameOne = "retry-1"
	var topicNameDlq = "retry.dlq"
	var topicNameDlqCustom = "mydlq"
	type args struct {
		msg   *kafka.Message
		retry config.Retry
	}
	tests := []struct {
		name    string
		args    args
		want    *kafka.Message
		wantErr bool
	}{
		{
			name: "withoutOrigin",
			args: args{
				msg: &kafka.Message{},
			},
			wantErr: true,
		},
		{
			name: "retry-1",
			args: args{
				retry: config.Retry{
					Attempts: []config.Attempt{
						{},
					},
				},
				msg: &kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic:     &topicName,
						Partition: -1,
					},
					Headers: []kafka.Header{
						{
							Key:   Origin,
							Value: []byte(topicName),
						},
					},
				},
			},
			want: &kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &topicNameOne,
					Partition: -1,
				},
				Timestamp:     time.Time{},
				TimestampType: 0,
				Headers: []kafka.Header{
					{
						Key:   Origin,
						Value: []byte(topicName),
					},
					{
						Key:   Attempt,
						Value: []byte("1"),
					},
				},
			},
		},
		{
			name: "dlq",
			args: args{
				retry: config.Retry{
					Dlq: "dlq",
					Attempts: []config.Attempt{
						{},
					},
				},
				msg: &kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic:     &topicName,
						Partition: -1,
					},
					Headers: []kafka.Header{
						{
							Key:   Origin,
							Value: []byte(topicName),
						},
						{
							Key:   Attempt,
							Value: []byte("1"),
						},
					},
				},
			},
			want: &kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &topicNameDlq,
					Partition: -1,
				},
				Timestamp:     time.Time{},
				TimestampType: 0,
				Headers: []kafka.Header{
					{
						Key:   Origin,
						Value: []byte(topicName),
					},
					{
						Key:   Attempt,
						Value: []byte("2"),
					},
				},
			},
		},
		{
			name: "dlq-with-header",
			args: args{
				retry: config.Retry{
					Dlq: "dlq",
					Attempts: []config.Attempt{
						{},
					},
				},
				msg: &kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic:     &topicName,
						Partition: -1,
					},
					Headers: []kafka.Header{
						{
							Key:   Origin,
							Value: []byte(topicName),
						},
						{
							Key:   Attempt,
							Value: []byte("1"),
						},
						{
							Key:   DlqTopic,
							Value: []byte("mydlq"),
						},
					},
				},
			},
			want: &kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &topicNameDlqCustom,
					Partition: -1,
				},
				Timestamp:     time.Time{},
				TimestampType: 0,
				Headers: []kafka.Header{
					{
						Key:   Origin,
						Value: []byte(topicName),
					},
					{
						Key:   Attempt,
						Value: []byte("2"),
					},
					{
						Key:   DlqTopic,
						Value: []byte("mydlq"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRetryMessage(tt.args.msg, tt.args.retry)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRetryMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRetryMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_copyMessage(t *testing.T) {
	type args struct {
		msg *kafka.Message
	}
	tests := []struct {
		name string
		args args
		want kafka.Message
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := copyMessage(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("copyMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAttempts(t *testing.T) {
	type args struct {
		msg *kafka.Message
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAttempts(tt.args.msg); got != tt.want {
				t.Errorf("getAttempts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTopicName(t *testing.T) {
	type args struct {
		msg      *kafka.Message
		retry    config.Retry
		attempts int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTopicName(tt.args.msg, tt.args.retry, tt.args.attempts); got != tt.want {
				t.Errorf("getTopicName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newMessage(t *testing.T) {
	type args struct {
		msg      *kafka.Message
		topic    string
		attempts int
	}
	tests := []struct {
		name    string
		args    args
		want    *kafka.Message
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newMessage(tt.args.msg, tt.args.topic, tt.args.attempts)
			if (err != nil) != tt.wantErr {
				t.Errorf("newMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}
