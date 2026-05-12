package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName              string
	JWTSecret            string
	AccessTokenExpireIn  string
	RefreshTokenExpireIn string
}

var Cfg Config

func LoadConfig() {
	godotenv.Load("./configs/.env")
	Cfg = Config{
		AppName:              os.Getenv("APP_NAME"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		AccessTokenExpireIn:  os.Getenv("JWT_EXPIRES_IN"),
		RefreshTokenExpireIn: os.Getenv("JWT_RefreshToken_EXPIRES_IN"),
	}
}
