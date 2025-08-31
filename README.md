# L0 Order Service

Микросервис для обработки заказов с использованием Kafka, PostgreSQL и кэширования в памяти.

## 🚀 Возможности

- **Kafka Consumer**: Получение заказов из очереди Kafka
- **PostgreSQL**: Надежное хранение данных с транзакциями  
- **In-Memory Cache**: Быстрый доступ к данным через кэш
- **REST API**: HTTP API для получения заказов по ID
- **Веб-интерфейс**: Простая веб-страница для поиска заказов
- **Docker**: Полная контейнеризация всех компонентов
- **Kafka Producer**: Скрипт-эмулятор для генерации тестовых заказов

## 📋 Архитектура

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Kafka Producer │───▶│     Kafka       │───▶│  Order Service  │
│   (Generator)   │    │   (Message      │    │                 │
└─────────────────┘    │    Broker)      │    │  ┌──────────────┤
                       └─────────────────┘    │  │ HTTP Server  │
┌─────────────────┐                           │  │ (Port 8080)  │
│   Web Browser   │◀──────────────────────────┤  └──────────────┤
└─────────────────┘                           │  ┌──────────────┤
                                              │  │    Cache     │
┌─────────────────┐                           │  │ (In-Memory)  │
│   PostgreSQL    │◀──────────────────────────┤  └──────────────┤
│  (Port 5433)    │                           └─────────────────┘
└─────────────────┘
```

## 🛠️ Быстрый запуск

### 1. Запуск всей системы
```bash
docker-compose up -d
```

### 2. Проверка состояния сервисов
```bash
docker-compose ps
```

### 3. Отправка тестовых заказов
```bash
# Переход в папку со скриптами
cd producer

# Отправка одного заказа
go run kafka_producer.go -count=1

# Отправка 5 заказов с задержкой 2 секунды
go run kafka_producer.go -count=5 -delay=2s

# Или использовать готовые скрипты
send_orders.bat          # Windows batch
demo.bat                 # Демонстрация системы
```

### 4. Использование веб-интерфейса
Откройте http://localhost:8080/ в браузере и введите Order UID для поиска.

## 🔧 Доступные сервисы

| Сервис | Адрес | Описание |
|--------|-------|----------|
| **Order Service** | http://localhost:8080/ | Веб-интерфейс для поиска заказов |
| **Order API** | http://localhost:8080/order/{order_uid} | REST API для получения заказа |
| **Kafka UI** | http://localhost:8081/ | Веб-интерфейс для управления Kafka |
| **PostgreSQL** | localhost:5433 | База данных (user: order_user, db: orders_service) |
| **Kafka** | localhost:9093 | Внешний порт для подключения к Kafka |

## 📊 Тестирование системы

### API тестирование
```bash
# Получение заказа по ID
curl http://localhost:8080/order/order_1_1234567890

# PowerShell
Invoke-WebRequest -Uri http://localhost:8080/order/order_1_1234567890 -UseBasicParsing
```

### Проверка логов
```bash
# Логи основного сервиса
docker-compose logs app

# Логи всех сервисов  
docker-compose logs

# Следить за логами в реальном времени
docker-compose logs -f app
```

### Мониторинг Kafka
1. Откройте http://localhost:8081/
2. Перейдите в Topics → orders
3. Просмотрите сообщения и консьюмеры

## 🎯 Структура проекта

```
L0/
├── cmd/server/           # Точка входа приложения
│   └── main.go
├── internal/             # Внутренняя логика
│   ├── cache/           # Кэширование в памяти  
│   ├── config/          # Конфигурация
│   ├── database/        # Работа с PostgreSQL
│   ├── handlers/        # HTTP обработчики
│   ├── kafka/          # Kafka consumer
│   ├── models/         # Модели данных
│   └── web/            # Веб-интерфейс
├── scripts/            # Скрипты для тестирования
│   ├── kafka_producer.go  # Генератор заказов
│   ├── Send-Orders.ps1    # PowerShell скрипт
│   ├── send_orders.bat    # Batch скрипт  
│   ├── demo.bat          # Демонстрация
│   └── README.md         # Документация скриптов
├── web/                # Веб-ресурсы
│   ├── static/         # CSS, JS файлы
│   └── templates/      # HTML шаблоны
├── migrations/         # SQL миграции
├── docker-compose.yml  # Docker конфигурация
└── Dockerfile         # Образ приложения
```

## Особенности реализации

### Производительность
- **Кэш в памяти** с concurrent-safe доступом (sync.RWMutex)
- **Connection pooling** для PostgreSQL
- **Batch операции** для массовой обработки

### Надежность  
- **Транзакции** в PostgreSQL для консистентности данных
- **Graceful shutdown** с корректным закрытием соединений
- **Error handling** на всех уровнях приложения
- **Валидация данных** входящих сообщений

### Масштабируемость
- **Stateless архитектура** для горизонтального масштабирования  
- **Kafka consumer groups** для распределения нагрузки
- **Docker контейнеры** для простого развертывания


## 🗄️ Схема базы данных

```sql
-- Основные таблицы
orders (основная таблица заказов)
├── order_uid TEXT PRIMARY KEY           -- Уникальный идентификатор заказа
├── track_number TEXT                    -- Номер отслеживания
├── entry TEXT                          -- Точка входа
├── delivery_id BIGINT → delivery(id)   -- FK к таблице доставки
├── payment_id BIGINT → payment(id)     -- FK к таблице платежей
├── locale TEXT                         -- Локализация
├── internal_signature TEXT             -- Внутренняя подпись
├── customer_id TEXT                    -- ID клиента
├── delivery_service TEXT               -- Служба доставки
├── shardkey TEXT                       -- Ключ шардирования
├── sm_id INTEGER                       -- ID сервиса
├── date_created TIMESTAMPTZ            -- Дата создания
└── oof_shard TEXT                      -- Шард

delivery (информация о доставке)
├── id BIGSERIAL PRIMARY KEY
├── name TEXT                           -- Имя получателя
├── phone TEXT                          -- Телефон
├── zip TEXT                            -- Почтовый код
├── city TEXT                           -- Город
├── address TEXT                        -- Адрес
├── region TEXT                         -- Регион
└── email TEXT                          -- Email

payment (информация о платеже)
├── id BIGSERIAL PRIMARY KEY
├── transaction TEXT                    -- ID транзакции
├── request_id TEXT                     -- ID запроса
├── currency TEXT                       -- Валюта
├── provider TEXT                       -- Провайдер платежа
├── amount BIGINT                       -- Сумма платежа
├── payment_dt BIGINT                   -- Дата платежа (timestamp)
├── bank TEXT                           -- Банк
├── delivery_cost BIGINT                -- Стоимость доставки
├── goods_total BIGINT                  -- Общая стоимость товаров
└── custom_fee BIGINT                   -- Таможенная пошлина

items (товары в заказе)
├── id BIGSERIAL PRIMARY KEY
├── chrt_id BIGINT                      -- ID характеристики
├── track_number TEXT                   -- Трек номер товара
├── price BIGINT                        -- Цена
├── rid TEXT                            -- RID
├── name TEXT                           -- Название товара
├── sale INTEGER                        -- Скидка (%)
├── size TEXT                           -- Размер
├── total_price BIGINT                  -- Итоговая цена
├── nm_id BIGINT                        -- Номенклатурный номер
├── brand TEXT                          -- Бренд
├── status INTEGER                      -- Статус товара
└── order_uid TEXT → orders(order_uid)  -- FK к заказу
```

### Индексы
- `idx_orders_customer_id` - по ID клиента
- `idx_items_order_uid` - по UID заказа
- `idx_orders_track_number` - по номеру отслеживания

### Связи
- `orders.delivery_id → delivery.id` (один к одному)
- `orders.payment_id → payment.id` (один к одному)
- `orders.order_uid ← items.order_uid` (один ко многим)

## 🔧 Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DB_HOST` | Хост PostgreSQL | localhost |
| `DB_PORT` | Порт PostgreSQL | 5432 |
| `DB_USER` | Пользователь БД | order_user |
| `DB_PASSWORD` | Пароль БД | order_pass |
| `DB_NAME` | Название БД | orders_service |
| `KAFKA_BROKER` | Адрес Kafka брокера | localhost:9092 |
| `SERVER_PORT` | Порт HTTP сервера | :8080 |

##  Примеры использования

### Отправка заказа через Kafka Producer
```bash
cd producer
go run kafka_producer.go -count=1 -delay=0s
```

### Получение заказа через API
```bash
curl http://localhost:8080/order/order_1_1234567890
```

### Использование веб-интерфейса
1. Перейдите на http://localhost:8080/
2. Введите Order UID в поле поиска
3. Нажмите "Найти"
4. Получите JSON с данными заказа

##  Полный тест системы

```bash
# 1. Запуск системы
docker-compose up -d

# 2. Подождать запуска (30-60 секунд)
docker-compose logs app

# 3. Отправить тестовые заказы  
cd producer
go run kafka_producer.go -count=3 -delay=1s

# 4. Проверить обработку в логах
docker-compose logs app | grep "успешно обработан"

# 5. Получить заказ через API
curl http://localhost:8080/order/order_1_XXXXXXXXX

# 6. Открыть веб-интерфейс
start http://localhost:8080/
```

## Остановка системы

```bash
# Остановка всех сервисов
docker-compose down

# Остановка с удалением томов (осторожно - удалит данные!)
docker-compose down -v
```

