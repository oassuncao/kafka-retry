package config

import (
	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"reflect"
)

func NewConsumer(group string) Consumer {
	consumer := Consumer{
		GroupId:         group,
		AutoOffsetReset: Earliest,
		MaxPoolInternal: 5400000, // 1h30m
	}
	return consumer
}

func LoadConsumer(kafka Kafka, consumer Consumer) confluent.ConfigMap {
	configMap := confluent.ConfigMap{}
	load(kafka, &configMap)
	load(consumer, &configMap)
	return configMap
}

func LoadProducer(kafka Kafka) confluent.ConfigMap {
	configMap := confluent.ConfigMap{}
	load(kafka, &configMap)
	return configMap
}

func load(config interface{}, configMap *confluent.ConfigMap) {
	valueOf := reflect.ValueOf(config)
	switch valueOf.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		load(valueOf.Elem(), configMap)
	case reflect.Struct:
		typeOf := valueOf.Type()
		for i := 0; i < valueOf.NumField(); i++ {
			fieldValue := valueOf.Field(i)
			fieldType := typeOf.Field(i)
			if fieldValue.Kind() == reflect.Struct {
				load(fieldValue.Interface(), configMap)
				continue
			}

			if value, ok := fieldType.Tag.Lookup("kafka"); ok {
				if fieldValue.Interface() != reflect.Zero(fieldType.Type).Interface() {
					configMap.SetKey(value, fieldValue.Interface())
				}
			}
		}
	}
}
