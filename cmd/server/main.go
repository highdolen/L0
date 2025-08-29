package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Создаём контекст для graceful shutdown
	ctxWithCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	// Запускаем Kafka Consumer в горутине
	go consumer.Start(ctxWithCancel)

	// Подключаем handlers
	r := mux.NewRouter()
	orderHandler := handlers.NewOrderHandler(orderCache, repo)
	r.HandleFunc("/order/{order_uid}", orderHandler.GetOrder).Methods("GET", "OPTIONS")

	// Подключаем веб-интерфейс
	web.RegisterWebHandlers(r)

	// Middleware
	r.Use(handlers.LoggingMiddleware)
	r.Use(handlers.CORSMiddleware)

	// Создаём HTTP сервер
	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}

	// Канал для перехвата сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем HTTP сервер в горутине
	go func() {
		log.Printf("HTTP сервер запущен на %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	<-sigChan
	log.Println("Получен сигнал завершения, начинаем graceful shutdown...")

	// Отменяем контекст для остановки Kafka consumer
	cancel()

	// Создаём контекст с таймаутом для shutdown HTTP сервера
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Graceful shutdown HTTP сервера
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Ошибка при graceful shutdown HTTP сервера: %v", err)
	} else {
		log.Println("HTTP сервер успешно остановлен")
	}

	// Закрываем Kafka consumer
	consumer.Close()
	log.Println("Kafka consumer успешно остановлен")

	log.Println("Graceful shutdown завершён")
}
