package config

type Consumer struct {
	GroupId         string `mapstructure:"group-id" kafka:"group.id"`
	AutoOffsetReset string `mapstructure:"auto-offset-reset" kafka:"auto.offset.reset"`
	MaxPoolInternal int    `mapstructure:"max-pool-internal" kafka:"max.poll.interval.ms"`
}

type AutoOffsetReset string

const (
	Earliest string = "earliest"
	Latest   string = "latest"
)
