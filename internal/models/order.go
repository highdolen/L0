package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Delivery struct {
	ID      int64  `json:"-" db:"id"`
	Name    string `json:"name" db:"name" validate:"required"`
	Phone   string `json:"phone" db:"phone" validate:"required"`
	Zip     string `json:"zip" db:"zip" validate:"required"`
	City    string `json:"city" db:"city" validate:"required"`
	Address string `json:"address" db:"address" validate:"required"`
	Region  string `json:"region" db:"region" validate:"required"`
	Email   string `json:"email" db:"email" validate:"required,email"`
}

type Payment struct {
	ID           int64  `json:"-" db:"id"`
	Transaction  string `json:"transaction" db:"transaction" validate:"required"`
	RequestID    string `json:"request_id" db:"request_id"`
	Currency     string `json:"currency" db:"currency" validate:"required"`
	Provider     string `json:"provider" db:"provider" validate:"required"`
	Amount       int64  `json:"amount" db:"amount" validate:"gte=0"`
	PaymentDt    int64  `json:"payment_dt" db:"payment_dt" validate:"gte=0"`
	Bank         string `json:"bank" db:"bank" validate:"required"`
	DeliveryCost int64  `json:"delivery_cost" db:"delivery_cost" validate:"gte=0"`
	GoodsTotal   int64  `json:"goods_total" db:"goods_total" validate:"gte=0"`
	CustomFee    int64  `json:"custom_fee" db:"custom_fee" validate:"gte=0"`
}

type Item struct {
	ID          int64  `json:"-" db:"id"`
	ChrtID      int64  `json:"chrt_id" db:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" db:"track_number" validate:"required"`
	Price       int64  `json:"price" db:"price" validate:"gte=0"`
	Rid         string `json:"rid" db:"rid" validate:"required"`
	Name        string `json:"name" db:"name" validate:"required"`
	Sale        int    `json:"sale" db:"sale" validate:"gte=0"`
	Size        string `json:"size" db:"size" validate:"required"`
	TotalPrice  int64  `json:"total_price" db:"total_price" validate:"gte=0"`
	NmID        int64  `json:"nm_id" db:"nm_id" validate:"required"`
	Brand       string `json:"brand" db:"brand" validate:"required"`
	Status      int    `json:"status" db:"status" validate:"gte=0"`
	OrderUID    string `json:"-" db:"order_uid"`
}

type Order struct {
	OrderUID          string    `json:"order_uid" db:"order_uid" validate:"required"`
	TrackNumber       string    `json:"track_number" db:"track_number" validate:"required"`
	Entry             string    `json:"entry" db:"entry" validate:"required"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Items             []Item    `json:"items" validate:"required,min=1,dive"`
	Locale            string    `json:"locale" db:"locale" validate:"required"`
	InternalSignature string    `json:"internal_signature" db:"internal_signature"`
	CustomerID        string    `json:"customer_id" db:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service" db:"delivery_service" validate:"required"`
	Shardkey          string    `json:"shardkey" db:"shardkey" validate:"required"`
	SmID              int       `json:"sm_id" db:"sm_id" validate:"gte=0"`
	DateCreated       time.Time `json:"date_created" db:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard" db:"oof_shard" validate:"required"`
}

// Validate — универсальная функция валидации для Order
func (o *Order) Validate() error {
	validate := validator.New()
	return validate.Struct(o)
}
