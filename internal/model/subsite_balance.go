package model

import "time"

// SubsiteBalance 分站余额
type SubsiteBalance struct {
	ID        int64     `json:"id" db:"id"`
	OwnerID   int64     `json:"owner_id" db:"owner_id"`
	Amount    float64   `json:"amount" db:"amount"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// SubsiteBalanceLog 分站余额变动记录
type SubsiteBalanceLog struct {
	ID            int64     `json:"id" db:"id"`
	OwnerID       int64     `json:"owner_id" db:"owner_id"`
	OrderID       int64     `json:"order_id,omitempty" db:"order_id"`
	Amount        float64   `json:"amount" db:"amount"`
	BeforeBalance float64   `json:"before_balance" db:"before_balance"`
	AfterBalance  float64   `json:"after_balance" db:"after_balance"`
	Type          int       `json:"type" db:"type"`
	Remark        string    `json:"remark" db:"remark"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
} 