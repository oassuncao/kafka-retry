package config

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"reflect"
	"testing"
)

func TestLoadConsumer(t *testing.T) {
	type args struct {
		kafka    Kafka
		consumer Consumer
	}
	tests := []struct {
		name string
		args args
		want kafka.ConfigMap
	}{
		{
			name: "default",
			args: args{
				kafka: Kafka{BootstrapServers: "localhost"},
				consumer: Consumer{
					GroupId:         "groupId",
					AutoOffsetReset: Earliest,
					MaxPoolInternal: 1000,
				},
			},
			want: kafka.ConfigMap{
				"bootstrap.servers":    "localhost",
				"group.id":             "groupId",
				"auto.offset.reset":    Earliest,
				"max.poll.interval.ms": 1000,
			},
		},
		{
			name: "wrong",
			args: args{
				kafka: Kafka{BootstrapServers: "newLocal"},
			},
			want: kafka.ConfigMap{
				"bootstrap.servers": "newLocal",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadConsumer(tt.args.kafka, tt.args.consumer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConsumer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadProducer(t *testing.T) {
	type args struct {
		kafka Kafka
	}
	tests := []struct {
		name string
		args args
		want kafka.ConfigMap
	}{
		{
			name: "default",
			args: args{
				kafka: Kafka{BootstrapServers: "localhost"},
			},
			want: kafka.ConfigMap{
				"bootstrap.servers": "localhost",
			},
		},
		{
			name: "wrong",
			args: args{
				kafka: Kafka{BootstrapServers: "newLocal"},
			},
			want: kafka.ConfigMap{
				"bootstrap.servers": "newLocal",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadProducer(tt.args.kafka); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadProducer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConsumer(t *testing.T) {
	type args struct {
		group string
	}
	tests := []struct {
		name string
		args args
		want Consumer
	}{
		{
			name: "default",
			args: args{
				group: "groupId",
			},
			want: Consumer{
				GroupId:         "groupId",
				AutoOffsetReset: Earliest,
				MaxPoolInternal: 5400000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConsumer(tt.args.group); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConsumer() = %v, want %v", got, tt.want)
			}
		})
	}
}
