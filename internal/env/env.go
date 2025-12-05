package env

import (
	"os"
	"strconv"
)

type Config struct {
	Port      int
	MongoURL  string
	RabbitURL string
	AuthURL   string
}

var config *Config

func Get() *Config {
	if config != nil {
		return config
	}

	portStr := getEnv("PORT", "3006")
	port, _ := strconv.Atoi(portStr)

	config = &Config{
		Port:      port,
		MongoURL:  getEnv("MONGO_URL", "mongodb://mongo:27017"),
		RabbitURL: getEnv("RABBIT_URL", "amqp://rabbitmq:5672"),
		AuthURL:   getEnv("AUTHGO_URL", "http://prod-auth-go:3000/v1/users/current"),
	}

	return config
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
