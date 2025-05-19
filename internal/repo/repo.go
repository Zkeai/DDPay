package repo

import (
	"context"
	"github.com/Zkeai/DDPay/internal/model"

	"github.com/Zkeai/DDPay/internal/conf"
	"github.com/Zkeai/DDPay/internal/repo/db"
)

type Repo struct {
	db *db.DB
}

func NewRepo(conf *conf.Conf) *Repo {
	return &Repo{
		db: db.NewDB(conf.DB),
	}
}

func (r *Repo) InsertMerchantWallet(ctx context.Context, req *model.MerchantWallet) error {

	return r.db.InsertMerchantWallet(ctx, req)
}

func (r *Repo) GetWalletByMerchantAndChain(ctx context.Context, merchantID uint32, chain string) (*model.MerchantWallet, error) {

	return r.db.GetWalletByMerchantAndChain(ctx, merchantID, chain)
}
