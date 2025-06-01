package model

import "time"

type OrderReq struct {
	Pid         uint32  `json:"pid"`
	TradeType   string  `json:"trade_type"`
	OrderId     string  `json:"order_id"`
	Amount      float64 `json:"amount"`
	Signature   string  `json:"signature"`
	NotifyUrl   string  `json:"notify_url"`
	RedirectUrl string  `json:"redirect_url"`
}

type OrderRsp struct {
	TradeId        string        `json:"trade_id"`
	OrderId        string        `json:"order_id"`
	Amount         float64       `json:"amount"`
	Token          string        `json:"token"`
	ActualAmount   float64       `json:"actual_amount"`
	ExpirationTime time.Duration `json:"expiration_time"`
	PaymentUrl     string        `json:"payment_url"`
}

type TradeOrders struct {
	Id          int64   `gorm:"primary_key;AUTO_INCREMENT;comment:id"`
	OrderId     string  `gorm:"column:order_id;type:varchar(128);not null;index;comment:商户ID"`
	TradeId     string  `gorm:"column:trade_id;type:varchar(128);not null;uniqueIndex;comment:本地ID"`
	TradeType   string  `gorm:"column:trade_type;type:varchar(20);not null;comment:交易类型"`
	TradeHash   string  `gorm:"column:trade_hash;type:varchar(128);default:'';unique;comment:交易哈希"`
	TradeRate   string  `gorm:"column:trade_rate;type:varchar(10);not null;comment:交易汇率"`
	Amount      string  `gorm:"type:decimal(10,2);not null;default:0;comment:交易数额"`
	Money       float64 `gorm:"type:decimal(10,2);not null;default:0;comment:订单交易金额"`
	Address     string  `gorm:"column:address;type:varchar(64);not null;comment:收款地址"`
	FromAddress string  `gorm:"type:varchar(34);not null;default:'';comment:支付地址"`
	Status      int     `gorm:"type:tinyint(1);not null;default:1;index;comment:交易状态"`
	Name        string  `gorm:"type:varchar(64);not null;default:'';comment:商品名称"`
	ApiType     string  `gorm:"type:varchar(20);not null;default:'epusdt';comment:API类型"`
	ReturnUrl   string  `gorm:"type:varchar(255);not null;default:'';comment:同步地址"`
	NotifyUrl   string  `gorm:"type:varchar(255);not null;default:'';comment:异步地址"`
	NotifyNum   int     `gorm:"column:notify_num;type:int(11);not null;default:0;comment:回调次数"`
	NotifyState int     `gorm:"column:notify_state;type:tinyint(1);not null;default:0;comment:回调状态 1：成功 0：失败"`
	RefBlockNum int64   `gorm:"type:bigint(20);not null;default:0;comment:交易所在区块"`
}
