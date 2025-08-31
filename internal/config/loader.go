package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

func LoadConfig() (*Config, error) {
    // Пытаемся загрузить .env файлы в порядке приоритета
    envFiles := []string{
        ".env",                    // Корень проекта (для Docker)
        "internal/config/.env",    // Локальная конфигурация
    }
    
    envLoaded := false
    for _, envFile := range envFiles {
        if err := godotenv.Load(envFile); err == nil {
            log.Printf("Загружен конфиг из %s", envFile)
            envLoaded = true
            break
        }
    }
    
    if !envLoaded {
        log.Println("No .env file found, using system environment variables")
    }

    cfg := &Config{
        DB: DBConfig{
            Host:     os.Getenv("DB_HOST"),
            Port:     os.Getenv("DB_PORT"),
            User:     os.Getenv("DB_USER"),
            Password: os.Getenv("DB_PASSWORD"),
            Name:     os.Getenv("DB_NAME"),
            SSLMode:  os.Getenv("DB_SSLMODE"),
        },
        Kafka: KafkaConfig{
            Broker: os.Getenv("KAFKA_BROKER"),
        },
        Server: ServerConfig{
            Port: os.Getenv("SERVER_PORT"),
        },
    }
    return cfg, nil
}
