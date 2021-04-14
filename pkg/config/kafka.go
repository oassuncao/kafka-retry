package config

type Kafka struct {
	BootstrapServers                   string `mapstructure:"bootstrap-servers" kafka:"bootstrap.servers"`
	SslCaLocation                      string `mapstructure:"ssl-ca-location" kafka:"ssl.ca.location"`
	SecurityProtocol                   string `mapstructure:"security-protocol" kafka:"security.protocol"`
	SaslUsername                       string `mapstructure:"sasl-username" kafka:"sasl.username"`
	SaslPassword                       string `mapstructure:"sasl-password" kafka:"sasl.password"`
	SaslMechanism                      string `mapstructure:"sasl-mechanism" kafka:"sasl.mechanism"`
	SslEndpointIdentificationAlgorithm string `mapstructure:"ssl-endpoint-identification-algorithm" kafka:"ssl.endpoint.identification.algorithm"`
}
