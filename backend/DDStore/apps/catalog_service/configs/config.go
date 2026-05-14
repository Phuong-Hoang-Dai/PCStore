package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName   string
	MongoURI  string
	DBName    string
	RedisAddr string
}

var Cfg Config

func LoadConfig() {
	godotenv.Load("./configs/.env")
	Cfg = Config{
		AppName:   os.Getenv("APP_NAME"),
		MongoURI:  os.Getenv("MONGO_URI"),
		DBName:    os.Getenv("DB_NAME"),
		RedisAddr: os.Getenv("REDIS_ADDR"),
	}
}