package model

import "time"

type MerchantWallet struct {
	ID             int64     `db:"id"`
	MerchantID     uint32    `db:"merchant_id"`
	Chain          string    `db:"chain"`           // "eth", "bsc", "solana","tron"
	Address        string    `db:"address"`         // 派生地址
	DerivationPath string    `db:"derivation_path"` // 派生路径
	CreatedAt      time.Time `db:"created_at"`
}

type RedisWallet struct {
	ID             int64
	MerchantID     uint32
	Chain          string
	Address        string
	DerivationPath string
	Amount         float64
	NotifyUrl      string
	RedirectUrl    string
	Offset         float64
	Status         string
}
