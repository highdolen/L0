package config

import "fmt"

type Config struct {
	DB     DBConfig
	Kafka  KafkaConfig
	Server ServerConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type KafkaConfig struct {
	Broker string
}

type ServerConfig struct {
	Port string
}

// Validate проверяет, что все обязательные поля заполнены
func (c *Config) Validate() error {
	if c.DB.Host == "" || c.DB.Port == "" || c.DB.User == "" || c.DB.Password == "" || c.DB.Name == "" {
		return fmt.Errorf("database configuration is incomplete")
	}
	if c.Kafka.Broker == "" {
		return fmt.Errorf("kafka broker address is missing")
	}
	if c.Server.Port == "" {
		return fmt.Errorf("server port is missing")
	}
	return nil
}
