package message

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"reflect"
	"testing"
)

func TestAddOrChange(t *testing.T) {
	message := &kafka.Message{
		Headers: []kafka.Header{
			{
				Key:   "header",
				Value: []byte("value"),
			},
		},
	}
	type args struct {
		msg   *kafka.Message
		name  string
		value []byte
	}

	type wants struct {
		returnValue bool
		value       []byte
	}
	tests := []struct {
		name string
		args args
		want wants
	}{
		{
			name: "add",
			args: args{
				msg:   message,
				name:  "newHeader",
				value: []byte("newHeader"),
			},
			want: wants{
				returnValue: true,
				value:       []byte("newHeader"),
			},
		},
		{
			name: "change",
			args: args{
				msg:   message,
				name:  "newHeader",
				value: []byte("changeHeader"),
			},
			want: wants{
				returnValue: false,
				value:       []byte("changeHeader"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddOrChange(tt.args.msg, tt.args.name, tt.args.value); got != tt.want.returnValue && reflect.DeepEqual(tt.args.msg.Headers[0].Value, tt.want.value) {
				t.Errorf("AddOrChange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHeader(t *testing.T) {
	type args struct {
		headers *[]kafka.Header
		name    string
	}
	tests := []struct {
		name string
		args args
		want *kafka.Header
	}{
		{
			name: "withHeader",
			args: args{
				headers: &([]kafka.Header{
					{
						Key:   "header",
						Value: []byte("value"),
					},
				}),
				name: "header",
			},
			want: &kafka.Header{
				Key:   "header",
				Value: []byte("value"),
			},
		},
		{
			name: "withOutHeader",
			args: args{
				headers: &([]kafka.Header{
					{
						Key:   "header",
						Value: []byte("value"),
					},
				}),
				name: "value",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHeader(tt.args.headers, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
