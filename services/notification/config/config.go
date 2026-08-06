package config

import (
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

var (
	cfg  Config
	once sync.Once
)

type Config struct {
	// Logging
	LogLevel string `envconfig:"LOG_LEVEL" default:"error"`

	// Ports
	HttpPort int `envconfig:"NOTIFICATION_HTTP_PORT" default:"8082"`
	GrpcPort int `envconfig:"GRPC_PORT" default:"50051"`

	// Kafka
	KafkaTopic string `envconfig:"KAFKA_TOPIC" required:"true"`
	KafkaHost  string `envconfig:"KAFKA_BROKER_HOST" required:"true"`
	KafkaPort  int    `envconfig:"KAFKA_BROKER_PORT" required:"true"`
}

func Load() *Config {
	once.Do(func() {
		var err error
		err = envconfig.Process("", &cfg)
		if err != nil {
			log.Fatalf("config error: %v", err)
		}
	})
	return &cfg
}
