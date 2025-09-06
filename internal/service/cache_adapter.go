package service

import (
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

// GetStats возвращает статистику кеша
func (a *cacheAdapter) GetStats() map[string]interface{} {
	return a.cache.GetStats()
}

// Close корректно завершает работу кеша
func (a *cacheAdapter) Close() {
	a.cache.Close()
}
