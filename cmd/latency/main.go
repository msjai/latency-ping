package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/msjai/latency-ping/internal/app"
	"github.com/msjai/latency-ping/internal/config"
)

func main() {
	// Устанавливаем переменные окружения из .env файла
	if err := godotenv.Load("./internal/config/.env"); err != nil {
		log.Print("No .env file found")
	}

	// Инициализируем конфиг
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
