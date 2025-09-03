package cache

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/highdolen/L0/internal/database"
	"github.com/highdolen/L0/internal/models"
)

// CacheEntry представляет запись в кеше с временной меткой
type CacheEntry struct {
	Order     models.Order
	Timestamp time.Time //момент, когда данные добавлены в кэш
}

type OrderCache struct {
	mu       sync.RWMutex
	cache    map[string]CacheEntry
	ttl      time.Duration
	stopChan chan bool
}

// New создает новый кэш с указанным TTL
func New(ttl time.Duration) *OrderCache {
	c := &OrderCache{
		cache:    make(map[string]CacheEntry),
		ttl:      ttl,
		stopChan: make(chan bool),
	}

	// Запускаем горутину для очистки устаревших записей
	go c.cleanupExpired()

	return c
}

// Get — получить заказ по UID с проверкой TTL
func (c *OrderCache) Get(uid string) (models.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.cache[uid]
	if !exists {
		return models.Order{}, false
	}

	// Проверяем, не истек ли TTL
	if time.Since(entry.Timestamp) > c.ttl {
		// Удаляем устаревшую запись
		delete(c.cache, uid)
		return models.Order{}, false
	}

	return entry.Order, true
}

// Set — добавить или обновить заказ с текущей временной меткой
func (c *OrderCache) Set(uid string, order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[uid] = CacheEntry{
		Order:     order,
		Timestamp: time.Now(),
	}
}

// Delete — удалить заказ по UID
func (c *OrderCache) Delete(uid string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, uid)
}

// LoadFromDB — восстановить кэш из базы
func (c *OrderCache) LoadFromDB(ctx context.Context, repo *database.OrderRepository) error {
	orders, err := repo.GetAllOrders(ctx)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for _, order := range orders {
		c.cache[order.OrderUID] = CacheEntry{
			Order:     order,
			Timestamp: now,
		}
	}

	log.Printf("Загружено %d заказов в кэш", len(orders))
	return nil
}

// Invalidate — инвалидировать конкретный заказ
func (c *OrderCache) Invalidate(uid string) {
	c.Delete(uid)
}

// InvalidateAll — очистить весь кэш
func (c *OrderCache) InvalidateAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]CacheEntry)
	log.Println("Кэш полностью очищен")
}

// Refresh — обновить конкретный заказ из базы данных
func (c *OrderCache) Refresh(ctx context.Context, uid string, repo *database.OrderRepository) error {
	order, err := repo.GetOrderByUID(ctx, uid)
	if err != nil {
		return err
	}

	if order != nil {
		c.Set(uid, *order)
		log.Printf("Заказ %s обновлен в кэше", uid)
	} else {
		c.Delete(uid)
		log.Printf("Заказ %s удален из кэша (не найден в БД)", uid)
	}

	return nil
}

// GetStats — получить статистику кэша
func (c *OrderCache) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	totalEntries := len(c.cache)
	expiredEntries := 0
	now := time.Now()

	for _, entry := range c.cache {
		if now.Sub(entry.Timestamp) > c.ttl {
			expiredEntries++
		}
	}

	return map[string]interface{}{
		"total_entries":   totalEntries,
		"expired_entries": expiredEntries,
		"valid_entries":   totalEntries - expiredEntries,
		"ttl_minutes":     c.ttl.Minutes(),
	}
}

// cleanupExpired — горутина для очистки устаревших записей
func (c *OrderCache) cleanupExpired() {
	ticker := time.NewTicker(c.ttl / 2) // На всякий случай проверяем каждые 15 минут, вдруг что то не удалилось
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()
			expiredKeys := make([]string, 0)

			for key, entry := range c.cache {
				if now.Sub(entry.Timestamp) >= c.ttl {
					expiredKeys = append(expiredKeys, key)
				}
			}

			for _, key := range expiredKeys {
				delete(c.cache, key)
			}

			if len(expiredKeys) > 0 {
				log.Printf("Удалено %d устаревших записей из кэша", len(expiredKeys))
			}
			c.mu.Unlock()

		case <-c.stopChan:
			return
		}
	}
}

// Close — остановить горутину очистки
func (c *OrderCache) Close() {
	close(c.stopChan)
}
