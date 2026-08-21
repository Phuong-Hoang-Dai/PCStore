package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName             string
	MongoURI            string
	DBName              string
	RedisAddr           string
	BrokerURI           string
	ProductCreatedQueue string
	ProductEditedQueue  string
	WorkerPerQueue      int
}

var Cfg Config

func LoadConfig() {
	godotenv.Load("./configs/.env")
	Cfg = Config{
		AppName:             os.Getenv("APP_NAME"),
		MongoURI:            os.Getenv("MONGO_URI"),
		DBName:              os.Getenv("DB_NAME"),
		RedisAddr:           os.Getenv("REDIS_ADDR"),
		BrokerURI:           getEnv("BROKER_URI", "amqp://guest:guest@localhost:5672/"),
		ProductCreatedQueue: getEnv("PRODUCT_CREATED_QUEUE", "ProductCreated"),
		ProductEditedQueue:  getEnv("PRODUCT_EDITED_QUEUE", "ProductEdited"),
		WorkerPerQueue:      getEnvInt("WORKER_PER_QUEUE", 20),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return value
}
