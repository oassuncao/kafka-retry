package config

import (
	"time"
)

type Application struct {
	Kafka Kafka
	Retry Retry
	Log   Log
}

type Retry struct {
	Topic           string
	Group           string
	Concurrency     int
	Dlq             string
	Attempts        []Attempt
	desiredAttempts *int
	Metric          Metric
}

type Attempt struct {
	Topic       string
	Delay       time.Duration
	Concurrency int
}

type Metric struct {
	Name        string
	Description string
}

type Log struct {
	Level       string
	Environment string
}

func (r Retry) GetDesiredRetries() int {
	return r.Concurrency
}

func (r Retry) GetDesiredAttempts() int {
	if r.desiredAttempts == nil {
		count := 0
		for _, attempt := range r.Attempts {
			count += attempt.Concurrency
		}
		r.desiredAttempts = &count
	}
	return *r.desiredAttempts
}
