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

// GetLoginLogs 获取登录日志（支持分页和筛选）
func (r *Repo) GetLoginLogs(ctx context.Context, params map[string]interface{}, page, pageSize int) ([]*model.LoginLog, int, error) {
	return r.db.GetLoginLogs(ctx, params, page, pageSize)
}

// ListUsers 获取用户列表（支持分页）
func (r *Repo) ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int, error) {
	return r.db.ListUsers(ctx, page, pageSize)
}

// 商户钱包相关方法
func (r *Repo) InsertMerchantWallet(ctx context.Context, req *model.MerchantWallet) error {
	return r.db.InsertMerchantWallet(ctx, req)
}

func (r *Repo) GetWalletByMerchantAndChain(ctx context.Context, merchantID uint32, chain string) (*model.MerchantWallet, error) {
	return r.db.GetWalletByMerchantAndChain(ctx, merchantID, chain)
}

// GetUserRepo 获取用户仓储接口
func (r *Repo) GetUserRepo() UserRepo {
	return NewUserRepo(r.db.GetDB())
}

// GetSubsiteRepo 获取分站仓储接口
func (r *Repo) GetSubsiteRepo() SubsiteRepo {
	return NewSubsiteRepo(r.db.GetDB())
}

// GetMembershipRepo 获取会员仓储接口
func (r *Repo) GetMembershipRepo() MembershipRepo {
	return NewMembershipRepo(r.db)
}

// 会员相关方法
func (r *Repo) GetMembershipLevels(ctx context.Context) ([]*model.MembershipLevel, error) {
	return r.db.GetMembershipLevels(ctx)
}

func (r *Repo) GetMembershipLevelByID(ctx context.Context, id int64) (*model.MembershipLevel, error) {
	return r.db.GetMembershipLevelByID(ctx, id)
}

func (r *Repo) GetMembershipLevelByLevel(ctx context.Context, level int) (*model.MembershipLevel, error) {
	return r.db.GetMembershipLevelByLevel(ctx, level)
}

func (r *Repo) GetMembershipBenefits(ctx context.Context, levelID int64) ([]*model.MembershipBenefit, error) {
	return r.db.GetMembershipBenefits(ctx, levelID)
}

func (r *Repo) GetMembershipRequirements(ctx context.Context, levelID int64) ([]*model.MembershipRequirement, error) {
	return r.db.GetMembershipRequirements(ctx, levelID)
}

func (r *Repo) GetUserMembership(ctx context.Context, userID int64) (*model.UserMembership, error) {
	return r.db.GetUserMembership(ctx, userID)
}

func (r *Repo) CreateUserMembership(ctx context.Context, membership *model.UserMembership) (int64, error) {
	return r.db.CreateUserMembership(ctx, membership)
}

func (r *Repo) UpdateUserMembership(ctx context.Context, membership *model.UserMembership) error {
	return r.db.UpdateUserMembership(ctx, membership)
}

func (r *Repo) ExtendUserMembership(ctx context.Context, userID int64, durationDays int) error {
	return r.db.ExtendUserMembership(ctx, userID, durationDays)
}

func (r *Repo) CreateMembershipTransaction(ctx context.Context, transaction *model.MembershipTransaction) (int64, error) {
	return r.db.CreateMembershipTransaction(ctx, transaction)
}

func (r *Repo) GetMembershipTransactions(ctx context.Context, userID int64, limit, offset int) ([]*model.MembershipTransaction, int, error) {
	return r.db.GetMembershipTransactions(ctx, userID, limit, offset)
}

func (r *Repo) UpdateMembershipTransaction(ctx context.Context, orderID string, status string) error {
	return r.db.UpdateMembershipTransaction(ctx, orderID, status)
}

func (r *Repo) GetMembershipTransactionByOrderID(ctx context.Context, orderID string) (*model.MembershipTransaction, error) {
	return r.db.GetMembershipTransactionByOrderID(ctx, orderID)
}
