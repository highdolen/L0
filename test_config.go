package main

import (
	"fmt"
	"log"

	"github.com/highdolen/L0/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	fmt.Printf("DB_HOST: %s\n", cfg.DB.Host)
	fmt.Printf("DB_USER: %s\n", cfg.DB.User)
	fmt.Printf("DB_PASSWORD: %s\n", cfg.DB.Password)
	fmt.Printf("KAFKA_BROKER: %s\n", cfg.Kafka.Broker)
	fmt.Printf("SERVER_PORT: %s\n", cfg.Server.Port)
}
