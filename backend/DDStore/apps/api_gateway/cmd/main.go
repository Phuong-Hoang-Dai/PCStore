package main

import (
	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/handler"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	configs.LoadConfig()
	handler.SetupHttp()
}
