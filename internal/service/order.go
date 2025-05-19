package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Zkeai/DDPay/internal/model"
	"github.com/Zkeai/DDPay/internal/wallet"
	"github.com/Zkeai/DDPay/pkg/redis"
	"sort"
	"time"
)

// CreateOrder 创建新的支付订单（含动态金额偏移逻辑）
func (s *Service) CreateOrder(ctx context.Context, req model.OrderReq) (model.OrderRsp, error) {
	// 查询商户钱包
	merchantWallet, err := s.repo.GetWalletByMerchantAndChain(ctx, req.Pid, req.TradeType)
	if err != nil {
		return model.OrderRsp{}, err
	}

	// 钱包不存在则派生一个新地址
	if merchantWallet == nil {
		deriver, err := wallet.GetGlobalDeriver()
		if err != nil {
			return model.OrderRsp{}, err
		}
		account, err := deriver.DeriveAccount(req.Pid)
		if err != nil {
			return model.OrderRsp{}, err
		}
		merchantWallet = &model.MerchantWallet{
			MerchantID:     req.Pid,
			Chain:          req.TradeType,
			Address:        account.Address,
			DerivationPath: fmt.Sprintf("m/44'/60'/0'/0/%d", req.Pid),
		}
		if err := s.repo.InsertMerchantWallet(ctx, merchantWallet); err != nil {
			return model.OrderRsp{}, err
		}
	}

	// 获取该地址当前未过期的订单（从 Redis 获取）
	pattern := fmt.Sprintf("%s_%d_*", req.TradeType, req.Pid)
	keys, err := redis.ScanKeys(pattern)
	if err != nil {
		return model.OrderRsp{}, fmt.Errorf("获取订单偏移失败: %w", err)
	}

	usedOffsets := make([]float64, 0)
	for _, k := range keys {
		val, err := redis.Get(k)
		if err != nil {
			continue
		}
		var cfg model.RedisWallet
		if err := json.Unmarshal([]byte(val), &cfg); err != nil {
			continue
		}
		if cfg.Address == merchantWallet.Address && cfg.Status == "pending" {
			usedOffsets = append(usedOffsets, cfg.Offset)
		}
	}

	// 计算最小可用 offset
	offset := getNextAmountOffset(usedOffsets)

	// 保存到 Redis
	cfg := model.RedisWallet{
		MerchantID:     req.Pid,
		Chain:          req.TradeType,
		Address:        merchantWallet.Address,
		Amount:         req.Amount, // 原始金额保持不变
		Offset:         offset,     // 偏移单独保存
		DerivationPath: merchantWallet.DerivationPath,
		NotifyUrl:      req.NotifyUrl,
		RedirectUrl:    req.RedirectUrl,
		Status:         "pending",
	}

	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return model.OrderRsp{}, err
	}

	key := fmt.Sprintf("%s_%d_%s", req.TradeType, req.Pid, req.OrderId)
	if err := redis.Set(key, jsonData, 30*time.Minute); err != nil {
		return model.OrderRsp{}, err
	}
	client := redis.GetClient()
	ttl, err := client.TTL(ctx, key).Result()

	return model.OrderRsp{
		TradeId:        key,
		OrderId:        req.OrderId,
		Amount:         req.Amount,
		ActualAmount:   req.Amount + offset,
		Token:          merchantWallet.Address,
		ExpirationTime: ttl,
		PaymentUrl:     "https://example.com/pay/checkout-counter/" + req.OrderId,
	}, nil
}

func (s *Service) GetOrderStatus(order string) model.RedisWallet {
	val, err := redis.Get(order)
	if err != nil {
		return model.RedisWallet{}
	}
	var cfg model.RedisWallet
	if err := json.Unmarshal([]byte(val), &cfg); err != nil {
		return model.RedisWallet{}
	}

	return cfg

}

// 计算最小未使用 offset（0.01步长）
func getNextAmountOffset(used []float64) float64 {
	sort.Float64s(used)
	step := 0.01
	for i := 0; ; i++ {
		next := float64(i+1) * step
		if i >= len(used) || used[i] > next {
			return next
		}
	}
}
