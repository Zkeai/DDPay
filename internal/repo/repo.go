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

// 用户相关方法
func (r *Repo) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	return r.db.GetUserByID(ctx, id)
}

func (r *Repo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return r.db.GetUserByEmail(ctx, email)
}

func (r *Repo) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	return r.db.CreateUser(ctx, user)
}

func (r *Repo) UpdateUser(ctx context.Context, user *model.User) error {
	return r.db.UpdateUser(ctx, user)
}

func (r *Repo) UpdatePassword(ctx context.Context, userID int64, password string) error {
	return r.db.UpdatePassword(ctx, userID, password)
}

func (r *Repo) UpdateLastLogin(ctx context.Context, userID int64, ip string) error {
	return r.db.UpdateLastLogin(ctx, userID, ip)
}

func (r *Repo) CreateVerificationCode(ctx context.Context, code *model.VerificationCode) error {
	return r.db.CreateVerificationCode(ctx, code)
}

func (r *Repo) GetVerificationCode(ctx context.Context, email, codeType string) (*model.VerificationCode, error) {
	return r.db.GetVerificationCode(ctx, email, codeType)
}

func (r *Repo) MarkVerificationCodeAsUsed(ctx context.Context, id int64) error {
	return r.db.MarkVerificationCodeAsUsed(ctx, id)
}

func (r *Repo) CreateOAuthAccount(ctx context.Context, account *model.OAuthAccount) error {
	return r.db.CreateOAuthAccount(ctx, account)
}

func (r *Repo) GetOAuthAccount(ctx context.Context, provider, providerUserID string) (*model.OAuthAccount, error) {
	return r.db.GetOAuthAccount(ctx, provider, providerUserID)
}

func (r *Repo) UpdateOAuthAccount(ctx context.Context, account *model.OAuthAccount) error {
	return r.db.UpdateOAuthAccount(ctx, account)
}

func (r *Repo) CreateLoginLog(ctx context.Context, log *model.LoginLog) error {
	return r.db.CreateLoginLog(ctx, log)
}

// 商户钱包相关方法
func (r *Repo) InsertMerchantWallet(ctx context.Context, req *model.MerchantWallet) error {
	return r.db.InsertMerchantWallet(ctx, req)
}

func (r *Repo) GetWalletByMerchantAndChain(ctx context.Context, merchantID uint32, chain string) (*model.MerchantWallet, error) {
	return r.db.GetWalletByMerchantAndChain(ctx, merchantID, chain)
}
