package service

import (
	"context"
	"log"

	"github.com/highdolen/L0/internal/cache"
	"github.com/highdolen/L0/internal/models"
)

// cacheAdapter адаптирует существующий кеш к интерфейсу CacheService
type cacheAdapter struct {
	cache *cache.OrderCache
}

// NewCacheAdapter создает новый адаптер для кеша
func NewCacheAdapter(cache *cache.OrderCache) CacheService {
	return &cacheAdapter{
		cache: cache,
	}
}

// Get получает заказ из кеша
func (a *cacheAdapter) Get(uid string) (models.Order, bool) {
	return a.cache.Get(uid)
}

// Set сохраняет заказ в кеш
func (a *cacheAdapter) Set(uid string, order models.Order) {
	a.cache.Set(uid, order)
}

// Delete удаляет заказ из кеша
func (a *cacheAdapter) Delete(uid string) {
	a.cache.Delete(uid)
}

// InvalidateAll очищает весь кеш
func (a *cacheAdapter) InvalidateAll() {
	a.cache.InvalidateAll()
}

// Refresh обновляет заказ в кеше из базы данных
func (a *cacheAdapter) Refresh(ctx context.Context, uid string, repo OrderRepository) error {
	// Получаем заказ из репозитория напрямую
	order, err := repo.GetOrderByUID(ctx, uid)
	if err != nil {
		return err
	}

	if order != nil {
		a.Set(uid, *order)
		log.Printf("Заказ %s обновлен в кэше", uid)
	} else {
		a.Delete(uid)
		log.Printf("Заказ %s удален из кэша (не найден в БД)", uid)
	}

	return nil
}

// GetStats возвращает статистику кеша
func (a *cacheAdapter) GetStats() map[string]interface{} {
	return a.cache.GetStats()
}

// LoadFromDB загружает данные из БД в кеш при инициализации
func (a *cacheAdapter) LoadFromDB(ctx context.Context, repo OrderRepository) error {
	// Получаем все заказы из репозитория
	orders, err := repo.GetAllOrders(ctx)
	if err != nil {
		return err
	}

	// Загружаем каждый заказ в кеш
	for _, order := range orders {
		a.Set(order.OrderUID, order)
	}

	log.Printf("Загружено %d заказов в кэш", len(orders))
	return nil
}

// Close корректно завершает работу кеша
func (a *cacheAdapter) Close() {
	a.cache.Close()
}
