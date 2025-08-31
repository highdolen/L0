package service

import (
	"context"
	"errors"
	"log"
)

// orderService реализует интерфейс OrderService
type orderService struct {
	repo  OrderRepository
	cache CacheService
}

// NewOrderService создает новый экземпляр сервиса заказов
func NewOrderService(repo OrderRepository, cache CacheService) OrderService {
	return &orderService{
		repo:  repo,
		cache: cache,
	}
}

// GetOrderByUID получает заказ по UID с использованием кеша
func (s *orderService) GetOrderByUID(ctx context.Context, uid string) (*OrderResult, error) {
	// Сначала проверяем кеш
	if cachedOrder, found := s.cache.Get(uid); found {
		log.Printf("Заказ %s получен из кеша", uid)
		return &OrderResult{
			Order:     &cachedOrder,
			FromCache: true,
		}, nil
	}

	// Если в кеше нет, идем в базу данных
	order, err := s.repo.GetOrderByUID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errors.New("заказ не найден")
	}

	// Кешируем заказ для будущих запросов
	s.cache.Set(uid, *order)
	log.Printf("Заказ %s загружен из БД и кеширован", uid)

	return &OrderResult{
		Order:     order,
		FromCache: false,
	}, nil
}

// GetOrderByUIDWithRefresh принудительно обновляет заказ из БД и возвращает его
func (s *orderService) GetOrderByUIDWithRefresh(ctx context.Context, uid string) (*OrderResult, error) {
	// Принудительно обновляем кеш из БД
	if err := s.cache.Refresh(ctx, uid, s.repo); err != nil {
		log.Printf("Ошибка обновления кеша для заказа %s: %v", uid, err)
		// Продолжаем выполнение, попытаемся получить из БД напрямую
	}

	// Получаем заказ напрямую из БД (после refresh)
	order, err := s.repo.GetOrderByUID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errors.New("заказ не найден")
	}

	return &OrderResult{
		Order:     order,
		FromCache: false, // Всегда false, так как это refresh
	}, nil
}

// GetCacheStats возвращает статистику кеша
func (s *orderService) GetCacheStats() map[string]interface{} {
	return s.cache.GetStats()
}

// InvalidateCache инвалидирует конкретный заказ в кеше
func (s *orderService) InvalidateCache(uid string) error {
	s.cache.Delete(uid)
	log.Printf("Инвалидирован кеш для заказа %s", uid)
	return nil
}

// InvalidateAllCache полностью очищает кеш
func (s *orderService) InvalidateAllCache() error {
	s.cache.InvalidateAll()
	log.Println("Весь кеш инвалидирован")
	return nil
}
