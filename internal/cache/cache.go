package cache

import (
	"context"
	"sync"

	"github.com/highdolen/L0/internal/database"
	"github.com/highdolen/L0/internal/models"
)

type OrderCache struct {
	mu    sync.RWMutex
	cache map[string]models.Order
}

// New создает новый кэш
func New() *OrderCache {
	return &OrderCache{
		cache: make(map[string]models.Order),
	}
}

// Get — получить заказ по UID
func (c *OrderCache) Get(uid string) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, exists := c.cache[uid]
	return order, exists
}

// Set — добавить или обновить заказ
func (c *OrderCache) Set(uid string, order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[uid] = order
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
	for _, order := range orders {
		c.cache[order.OrderUID] = order
	}

	return nil
}
