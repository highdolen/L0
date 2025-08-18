package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/highdolen/L0/internal/cache"
	"github.com/highdolen/L0/internal/config"
	"github.com/highdolen/L0/internal/database"
	"github.com/highdolen/L0/internal/handlers"
	"github.com/highdolen/L0/internal/kafka"
	"github.com/highdolen/L0/internal/web"
)

func main() {
	ctx := context.Background()

	// Загружаем конфиг
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Ошибка валидации конфигурации: %v", err)
	}

	// Подключение к базе
	dsn := "postgres://" + cfg.DB.User + ":" + cfg.DB.Password + "@" + cfg.DB.Host + ":" + cfg.DB.Port + "/" + cfg.DB.Name + "?sslmode=" + cfg.DB.SSLMode
	db, err := database.ConnectDB(dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	repo := database.NewOrderRepository(db)

	// Создаём кэш
	orderCache := cache.New()

	// Загружаем данные из БД в кэш
	if err := orderCache.LoadFromDB(ctx, repo); err != nil {
		log.Fatalf("Ошибка загрузки кэша: %v", err)
	}
	log.Println("Кэш успешно загружен")

	// Создаём Kafka Consumer
	consumer := kafka.NewConsumer(
		[]string{cfg.Kafka.Broker},
		"orders",
		"group-1",
		repo,
		orderCache,
	)

	// Запускаем Kafka Consumer в горутине
	go consumer.Start(ctx)

	// Подключаем handlers
	r := mux.NewRouter()
	orderHandler := handlers.NewOrderHandler(orderCache, repo)
	r.HandleFunc("/order/{order_uid}", orderHandler.GetOrder).Methods("GET", "OPTIONS")

	// Подключаем веб-интерфейс
	web.RegisterWebHandlers(r)

	// Middleware
	r.Use(handlers.LoggingMiddleware)
	r.Use(handlers.CORSMiddleware)

	log.Printf("HTTP сервер запущен на %s", cfg.Server.Port)
	if err := http.ListenAndServe(cfg.Server.Port, r); err != nil {
		log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
	}
}
