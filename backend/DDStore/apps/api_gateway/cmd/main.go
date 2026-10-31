package main

import (
	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/handler"
	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/logger"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	logger.Setup("api_gateway")
	configs.LoadConfig()
	handler.SetupHttp()
}
