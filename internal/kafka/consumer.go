package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/highdolen/L0/internal/cache"
	"github.com/highdolen/L0/internal/database"
	"github.com/highdolen/L0/internal/models"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader   *kafka.Reader
	repo     *database.OrderRepository
	cache    *cache.OrderCache
	validate *validator.Validate
}

func NewConsumer(brokers []string, topic, groupID string, repo *database.OrderRepository, cache *cache.OrderCache) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.LastOffset,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		MaxWait:     time.Second,
	})

	return &Consumer{
		reader:   r,
		repo:     repo,
		cache:    cache,
		validate: validator.New(),
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("Kafka consumer started...")
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Ошибка чтения сообщения из Kafka: %v", err)
			continue
		}

		var order models.Order
		if err := json.Unmarshal(m.Value, &order); err != nil {
			log.Printf("Ошибка парсинга JSON: %v", err)
			continue
		}

		if err := c.validate.Struct(order); err != nil {
			log.Printf("Ошибка валидации данных: %v", err)
			continue
		}

		if err := c.repo.CreateOrder(ctx, &order); err != nil {
			log.Printf("Ошибка сохранения заказа в БД: %v", err)
			continue
		}

		c.cache.Set(order.OrderUID, order)
		log.Printf("Заказ %s успешно обработан", order.OrderUID)
	}
}

func (c *Consumer) Close() {
	c.reader.Close()
}
