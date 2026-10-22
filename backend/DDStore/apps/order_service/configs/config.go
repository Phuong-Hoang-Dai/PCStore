package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName         string
	ConnectStr      string
	StripeSecretKey  string
	StripeSuccessURL string
	StripeCancelURL  string
	StripeWebhookKey string
}

var Cfg Config

func LoadConfig() {
	godotenv.Load("./configs/.env")
	Cfg = Config{
		AppName:          os.Getenv("APP_NAME"),
		ConnectStr:       os.Getenv("DB_CONN_STR"),
		StripeSecretKey:  os.Getenv("STRIPE_SECRET_KEY"),
		StripeSuccessURL: os.Getenv("STRIPE_SUCCESS_URL"),
		StripeCancelURL:  os.Getenv("STRIPE_CANCEL_URL"),
		StripeWebhookKey: os.Getenv("STRIPE_WEBHOOK_SECRET"),
	}
}
