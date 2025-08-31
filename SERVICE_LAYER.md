# Сервисный слой (Service Layer)

В рамках рефакторинга проекта была добавлена архитектурная прослойка - сервисный слой, который изолирует бизнес-логику от HTTP handlers и обеспечивает лучшую организацию кода.

## Структура сервисного слоя

### Интерфейсы (`internal/service/interfaces.go`)

1. **OrderService** - основной интерфейс для бизнес-логики работы с заказами
2. **OrderRepository** - интерфейс для работы с базой данных
3. **CacheService** - интерфейс для работы с кешем
4. **OrderResult** - структура результата с метаданными о источнике данных

### Реализации

#### OrderService (`internal/service/order_service.go`)
- `GetOrderByUID()` - получение заказа с автоматическим кешированием
- `GetOrderByUIDWithRefresh()` - принудительное обновление из БД
- `GetCacheStats()` - статистика кеша
- `InvalidateCache()` - инвалидация конкретного заказа
- `InvalidateAllCache()` - полная очистка кеша

#### Адаптеры
- **RepositoryAdapter** (`internal/service/repository_adapter.go`) - адаптирует существующий `database.OrderRepository` к интерфейсу `OrderRepository`
- **CacheAdapter** (`internal/service/cache_adapter.go`) - адаптирует существующий `cache.OrderCache` к интерфейсу `CacheService`

## Преимущества архитектуры

1. **Разделение ответственности** - handlers отвечают только за HTTP, сервис - за бизнес-логику
2. **Тестируемость** - легко создавать мок-объекты для интерфейсов
3. **Переиспользование** - логика может использоваться разными handlers или другими сервисами
4. **Расширяемость** - легко добавлять новую функциональность
5. **Инверсия зависимостей** - зависимости инжектятся через интерфейсы

## Метаданные запросов

Введена структура `OrderResult`, которая содержит:
- `Order` - сам заказ
- `FromCache` - флаг, указывающий, был ли заказ получен из кеша

Это позволяет корректно устанавливать HTTP-заголовок `X-Cache`.

## Использование в handlers

```go
// Создание сервиса
repoAdapter := service.NewRepositoryAdapter(repo)
cacheAdapter := service.NewCacheAdapter(orderCache)
orderService := service.NewOrderService(repoAdapter, cacheAdapter)

// Создание handler
orderHandler := handlers.NewOrderHandler(orderService)
```

Handler теперь использует только сервисный слой, не обращаясь напрямую к кешу или репозиторию.

## Обратная совместимость

Существующие компоненты (кеш и репозиторий) остались без изменений. Адаптеры обеспечивают совместимость с новыми интерфейсами.
