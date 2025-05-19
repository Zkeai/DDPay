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
