package model

import "time"

// Subsite 分站模型
type Subsite struct {
	ID             int64     `json:"id" db:"id"`
	OwnerID        int64     `json:"owner_id" db:"owner_id"`
	Name           string    `json:"name" db:"name"`
	Domain         string    `json:"domain" db:"domain"`
	Subdomain      string    `json:"subdomain" db:"subdomain"`
	Logo           string    `json:"logo" db:"logo"`
	Description    string    `json:"description" db:"description"`
	Theme          string    `json:"theme" db:"theme"`
	Status         int       `json:"status" db:"status"`
	CommissionRate float64   `json:"commission_rate" db:"commission_rate"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// SubsiteConfig 分站配置模型（JSON格式）
type SubsiteConfig struct {
	ID        int64     `json:"id" db:"id"`
	SubsiteID int64     `json:"subsite_id" db:"subsite_id"`
	Config    string    `json:"config" db:"config"` // 存储JSON字符串
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// SubsiteProduct 分站商品
type SubsiteProduct struct {
	ID            int64     `json:"id" db:"id"`
	SubsiteID     int64     `json:"subsite_id" db:"subsite_id"`
	MainProductID int64     `json:"main_product_id" db:"main_product_id"`
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description" db:"description"`
	Price         float64   `json:"price" db:"price"`
	OriginalPrice float64   `json:"original_price" db:"original_price"`
	Stock         int       `json:"stock" db:"stock"`
	Status        int       `json:"status" db:"status"`
	Image         string    `json:"image" db:"image"`
	IsTimeLimited int       `json:"is_time_limited" db:"is_time_limited"`
	StartTime     time.Time `json:"start_time" db:"start_time"`
	EndTime       time.Time `json:"end_time" db:"end_time"`
	SortOrder     int       `json:"sort_order" db:"sort_order"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// SubsiteWithdrawal 分站提现
type SubsiteWithdrawal struct {
	ID          int64     `json:"id" db:"id"`
	OwnerID     int64     `json:"owner_id" db:"owner_id"`
	Amount      float64   `json:"amount" db:"amount"`
	Status      int       `json:"status" db:"status"`
	AccountType string    `json:"account_type" db:"account_type"`
	AccountName string    `json:"account_name" db:"account_name"`
	AccountNo   string    `json:"account_no" db:"account_no"`
	Remark      string    `json:"remark" db:"remark"`
	AdminRemark string    `json:"admin_remark" db:"admin_remark"`
	ProcessedAt time.Time `json:"processed_at,omitempty" db:"processed_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// SubsiteInfo 分站信息（包含关联的用户信息）
type SubsiteInfo struct {
	Subsite      *Subsite     `json:"subsite"`
	Owner        *UserProfile `json:"owner"`
	ProductCount int          `json:"product_count"`
	OrderCount   int          `json:"order_count"`
	Balance      float64      `json:"balance"`
}

// 常量定义
const (
	// 分站状态
	SubsiteStatusDisabled = 0
	SubsiteStatusEnabled  = 1

	// 分站商品状态
	SubsiteProductStatusOffline = 0
	SubsiteProductStatusOnline  = 1

	// 分站订单状态
	SubsiteOrderStatusPending   = 0
	SubsiteOrderStatusPaid      = 1
	SubsiteOrderStatusCompleted = 2
	SubsiteOrderStatusCancelled = 3

	// 分站余额变动类型
	SubsiteBalanceTypeCommission = 1
	SubsiteBalanceTypeWithdrawal = 2
	SubsiteBalanceTypeAdjustment = 3

	// 分站提现状态
	SubsiteWithdrawalStatusPending   = 0
	SubsiteWithdrawalStatusCompleted = 1
	SubsiteWithdrawalStatusRejected  = 2
)