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
	"github.com/highdolen/L0/internal/service"
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

	// Создаём кэш с TTL 30 минут
	orderCache := cache.New(30 * time.Minute)

	// Создаём адаптеры для сервисного слоя
	repoAdapter := service.NewRepositoryAdapter(repo)
	cacheAdapter := service.NewCacheAdapter(orderCache)

	// Загружаем данные из БД в кэш через адаптер
	if err := cacheAdapter.LoadFromDB(ctx, repoAdapter); err != nil {
		log.Fatalf("Ошибка загрузки кэша: %v", err)
	}
	log.Println("Кэш успешно загружен")

	// Создаём сервис заказов
	orderService := service.NewOrderService(repoAdapter, cacheAdapter)

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

	// Подключаем handlers
	r := mux.NewRouter()
	orderHandler := handlers.NewOrderHandler(orderService)

	// API для работы с заказами
	r.HandleFunc("/order/{order_uid}", orderHandler.GetOrder).Methods("GET", "OPTIONS")

	// API для управления кешом
	r.HandleFunc("/cache/stats", orderHandler.GetCacheStats).Methods("GET", "OPTIONS")
	r.HandleFunc("/cache/invalidate/{order_uid}", orderHandler.InvalidateCache).Methods("POST", "DELETE", "OPTIONS")
	r.HandleFunc("/cache/invalidate", orderHandler.InvalidateCache).Methods("POST", "DELETE", "OPTIONS")

	err = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		log.Printf("Route registered: %s Methods: %v", path, methods)
		return nil
	})
	if err != nil {
		log.Printf("Ошибка при обходе маршрутов: %v", err)
	}

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

	// Канал для уведомления о завершении shutdown
	shutdownComplete := make(chan bool, 1)

	// Горутина для обработки сигналов и graceful shutdown
	go func() {
		// Ожидаем сигнал
		sig := <-sigChan
		log.Printf("Получен сигнал %v, начинаем graceful shutdown...", sig)

		// Отменяем контекст для остановки Kafka consumer
		log.Println("Останавливаем Kafka consumer...")
		cancel()

		// Ждем немного, чтобы consumer успел обработать отмену контекста
		time.Sleep(100 * time.Millisecond)

		// Создаём контекст с таймаутом для shutdown HTTP сервера
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		// Graceful shutdown HTTP сервера
		log.Println("Останавливаем HTTP сервер...")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Ошибка при graceful shutdown HTTP сервера: %v", err)
		} else {
			log.Println("HTTP сервер успешно остановлен")
		}

		// Закрываем Kafka consumer
		log.Println("Закрываем Kafka consumer...")
		consumer.Close()
		log.Println("Kafka consumer успешно остановлен")

		// Закрываем кеш (останавливаем горутину очистки)
		log.Println("Останавливаем кеш...")
		orderCache.Close()
		log.Println("Кеш успешно остановлен")

		log.Println("Graceful shutdown завершён")
		shutdownComplete <- true
	}()

	// Запускаем Kafka Consumer в горутине
	go consumer.Start(ctxWithCancel)

	// Запускаем HTTP сервер в горутине
	go func() {
		log.Printf("HTTP сервер запущен на %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Ошибка HTTP сервера: %v", err)
			// Отправляем сигнал для shutdown, если сервер упал
			sigChan <- syscall.SIGTERM
		}
	}()

	log.Println("Сервер запущен. Нажмите Ctrl+C для остановки.")

	// Ожидаем завершения shutdown
	<-shutdownComplete
}
