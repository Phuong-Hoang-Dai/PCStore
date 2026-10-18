package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName    string
	ConnectStr string
}

var Cfg Config

func LoadConfig() {
	godotenv.Load("./configs/.env")
	Cfg = Config{
		AppName:    os.Getenv("APP_NAME"),
		ConnectStr: os.Getenv("DB_CONN_STR"),
	}
}
