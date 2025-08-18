package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

// Структуры данных (копия из models)
type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int64  `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int64  `json:"delivery_cost"`
	GoodsTotal   int64  `json:"goods_total"`
	CustomFee    int64  `json:"custom_fee"`
}

type Item struct {
	ChrtID      int64  `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int64  `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int64  `json:"total_price"`
	NmID        int64  `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

type Order struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

func main() {
	var (
		broker = flag.String("broker", "localhost:9093", "Kafka broker address")
		topic  = flag.String("topic", "orders", "Kafka topic")
		count  = flag.Int("count", 1, "Number of messages to send")
		delay  = flag.Duration("delay", 1*time.Second, "Delay between messages")
	)
	flag.Parse()

	log.Printf("Подключаемся к Kafka брокеру: %s", *broker)
	log.Printf("Топик: %s", *topic)
	log.Printf("Количество сообщений: %d", *count)
	log.Printf("Задержка между сообщениями: %v", *delay)

	// Создаем Kafka writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{*broker},
		Topic:   *topic,
	})
	defer writer.Close()

	ctx := context.Background()

	for i := 0; i < *count; i++ {
		order := generateRandomOrder(i + 1)
		
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Printf("Ошибка кодирования JSON: %v", err)
			continue
		}

		message := kafka.Message{
			Key:   []byte(order.OrderUID),
			Value: orderJSON,
		}

		err = writer.WriteMessages(ctx, message)
		if err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
			continue
		}

		log.Printf("✅ Отправлен заказ #%d: %s", i+1, order.OrderUID)
		
		if i < *count-1 {
			time.Sleep(*delay)
		}
	}

	log.Printf("🎉 Успешно отправлено %d сообщений!", *count)
}

func generateRandomOrder(num int) Order {
	rand.Seed(time.Now().UnixNano() + int64(num))
	
	orderUID := fmt.Sprintf("order_%d_%d", num, time.Now().Unix())
	trackNumber := fmt.Sprintf("TRACK%06d", rand.Intn(999999))
	
	// Случайные города и имена
	cities := []string{"Москва", "Санкт-Петербург", "Казань", "Екатеринбург", "Новосибирск"}
	names := []string{"Иван Иванов", "Петр Петров", "Анна Сидорова", "Мария Козлова", "Алексей Смирнов"}
	brands := []string{"Nike", "Adidas", "Puma", "Reebok", "New Balance"}
	
	city := cities[rand.Intn(len(cities))]
	name := names[rand.Intn(len(names))]
	
	return Order{
		OrderUID:    orderUID,
		TrackNumber: trackNumber,
		Entry:       "WBIL",
		Delivery: Delivery{
			Name:    name,
			Phone:   fmt.Sprintf("+7%010d", rand.Int63n(9999999999)),
			Zip:     strconv.Itoa(100000 + rand.Intn(599999)),
			City:    city,
			Address: fmt.Sprintf("ул. %s, д. %d", generateStreetName(), rand.Intn(100)+1),
			Region:  city + " область",
			Email:   fmt.Sprintf("user%d@example.com", rand.Intn(9999)),
		},
		Payment: Payment{
			Transaction:  fmt.Sprintf("txn_%d", rand.Int63()),
			RequestID:    fmt.Sprintf("req_%d", rand.Int63()),
			Currency:     "RUB",
			Provider:     "alfabank",
			Amount:       int64(rand.Intn(10000) + 1000),
			PaymentDt:    time.Now().Unix(),
			Bank:         "alfa",
			DeliveryCost: int64(rand.Intn(500) + 200),
			GoodsTotal:   int64(rand.Intn(8000) + 1000),
			CustomFee:    0,
		},
		Items:             generateItems(trackNumber, rand.Intn(3)+1, brands),
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        fmt.Sprintf("customer_%d", rand.Intn(9999)),
		DeliveryService:   "meest",
		Shardkey:          strconv.Itoa(rand.Intn(10)),
		SmID:              rand.Intn(100),
		DateCreated:       time.Now(),
		OofShard:          "1",
	}
}

func generateItems(trackNumber string, count int, brands []string) []Item {
	items := make([]Item, count)
	itemNames := []string{
		"Футболка", "Джинсы", "Кроссовки", "Куртка", "Рубашка",
		"Свитер", "Брюки", "Платье", "Кепка", "Носки",
	}
	sizes := []string{"XS", "S", "M", "L", "XL", "XXL"}

	for i := 0; i < count; i++ {
		price := int64(rand.Intn(5000) + 500)
		sale := rand.Intn(50)
		totalPrice := price * int64(100-sale) / 100

		items[i] = Item{
			ChrtID:      int64(rand.Intn(999999) + 100000),
			TrackNumber: trackNumber,
			Price:       price,
			Rid:         fmt.Sprintf("rid_%d", rand.Int63()),
			Name:        itemNames[rand.Intn(len(itemNames))],
			Sale:        sale,
			Size:        sizes[rand.Intn(len(sizes))],
			TotalPrice:  totalPrice,
			NmID:        int64(rand.Intn(9999999) + 1000000),
			Brand:       brands[rand.Intn(len(brands))],
			Status:      202,
		}
	}

	return items
}

func generateStreetName() string {
	streets := []string{
		"Ленина", "Пушкина", "Гагарина", "Советская", "Мира",
		"Победы", "Молодежная", "Центральная", "Садовая", "Школьная",
	}
	return streets[rand.Intn(len(streets))]
}
