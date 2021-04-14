package config

import (
	"reflect"
	"testing"
)

func TestSetContext(t *testing.T) {
	basicConf := Application{
		Kafka: Kafka{
			BootstrapServers: "Host",
		},
		Retry: Retry{
			Topic:       "topic",
			Group:       "group",
			Concurrency: 1,
			Dlq:         "dlq",
			Attempts:    nil,
		},
	}

	type args struct {
		conf Application
	}
	tests := []struct {
		name string
		args args
		want Context
	}{
		{
			name: "default",
			args: args{
				conf: basicConf,
			},
			want: Context{
				Config: basicConf,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetContext(tt.args.conf)
			if got := GetContext(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
