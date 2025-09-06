package service

import (
	"context"

	"github.com/highdolen/L0/internal/models"
)

// OrderResult содержит результат получения заказа с метаданными
type OrderResult struct {
	Order     *models.Order
	FromCache bool
}

// OrderService определяет интерфейс для бизнес-логики работы с заказами
type OrderService interface {
	// Загружает все заказы из БД в кэш при инициализации
	LoadFromDB(ctx context.Context) error

	// Refresh обновляет конкретный заказ в кэше из базы данных
	Refresh(ctx context.Context, uid string) error

	// GetOrderByUID получает заказ по UID с использованием кеша
	GetOrderByUID(ctx context.Context, uid string) (*OrderResult, error)

	// GetOrderByUIDWithRefresh принудительно обновляет заказ из БД и возвращает его
	GetOrderByUIDWithRefresh(ctx context.Context, uid string) (*OrderResult, error)

	// GetCacheStats возвращает статистику кеша
	GetCacheStats() map[string]interface{}

	// InvalidateCache инвалидирует конкретный заказ в кеше
	InvalidateCache(uid string) error

	// InvalidateAllCache полностью очищает кеш
	InvalidateAllCache() error
}

// OrderRepository определяет интерфейс для работы с базой данных
type OrderRepository interface {
	// GetOrderByUID получает заказ из базы данных по UID
	GetOrderByUID(ctx context.Context, uid string) (*models.Order, error)

	// CreateOrder создает новый заказ в базе данных
	CreateOrder(ctx context.Context, order *models.Order) error

	// GetAllOrders получает все заказы из базы данных
	GetAllOrders(ctx context.Context) ([]models.Order, error)
}

// CacheService определяет интерфейс для работы с кешем
type CacheService interface {
	// Get получает заказ из кеша
	Get(uid string) (models.Order, bool)

	// Set сохраняет заказ в кеш
	Set(uid string, order models.Order)

	// Delete удаляет заказ из кеша
	Delete(uid string)

	// InvalidateAll очищает весь кеш
	InvalidateAll()

	// GetStats возвращает статистику кеша
	GetStats() map[string]interface{}

	// Close корректно завершает работу кеша
	Close()
}
