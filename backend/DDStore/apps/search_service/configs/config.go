package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName    string
	Port       string
	ESAddrs    string // comma-separated Elasticsearch node URLs
	ESUsername string
	ESPassword string
	ESIndex    string
}

var Cfg Config

func LoadConfig() {
	// .env is optional: in containers config comes from real environment vars.
	godotenv.Load("./configs/.env")

	Cfg = Config{
		AppName:    getEnv("APP_NAME", "PC Store Search"),
		Port:       getEnv("PORT", "8889"),
		ESAddrs:    getEnv("ES_ADDRESSES", "http://localhost:9200"),
		ESUsername: os.Getenv("ES_USERNAME"),
		ESPassword: os.Getenv("ES_PASSWORD"),
		ESIndex:    getEnv("ES_INDEX", "products"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
