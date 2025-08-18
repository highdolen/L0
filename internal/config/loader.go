package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

func LoadConfig() (*Config, error) {
    // Загружаем .env
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system env variables")
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
