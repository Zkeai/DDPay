package model

import "time"

// SubsiteOrder 分站订单
type SubsiteOrder struct {
	ID           int64     `json:"id" db:"id"`
	OrderNo      string    `json:"order_no" db:"order_no"`
	SubsiteID    int64     `json:"subsite_id" db:"subsite_id"`
	UserID       int64     `json:"user_id" db:"user_id"`
	ProductID    int64     `json:"product_id" db:"product_id"`
	Quantity     int       `json:"quantity" db:"quantity"`
	Amount       float64   `json:"amount" db:"amount"`
	Commission   float64   `json:"commission" db:"commission"`
	Status       int       `json:"status" db:"status"`
	PayTime      time.Time `json:"pay_time,omitempty" db:"pay_time"`
	CompleteTime time.Time `json:"complete_time,omitempty" db:"complete_time"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
} 