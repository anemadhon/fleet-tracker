package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string
	LogDir string

	PostgresDSN string
	RedisAddr   string

	MQTTBroker string
	MQTTPort   string
	MQTTUseTLS string

	RabbitURL string
}

var Cfg *Config

func Load() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Failed to load env")
	}

	Cfg = &Config{
		AppEnv: getEnv("APP_ENV", "development"),
		LogDir: getEnv("LOG_DIR", "./logs"),

		PostgresDSN: getEnv("POSTGRES_DSN", "postgres://postgres:postgrePWD.90@localhost:5432/rnd?sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),

		MQTTBroker: getEnv("MQTT_BROKER", "localhost"),
		MQTTPort:   getEnv("MQTT_PORT", "1883"),
		MQTTUseTLS: getEnv("MQTT_USE_TLS", "false"),

		RabbitURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5673/"),
	}

	log.Printf("config loaded: ENV=%s", Cfg.AppEnv)
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return def
}
