package repo

import (
	"context"

	"github.com/Zkeai/DDPay/internal/model"
	"github.com/Zkeai/DDPay/internal/repo/db"
)

// MembershipRepo 会员等级存储库接口
type MembershipRepo interface {
	// 会员等级相关
	GetMembershipLevels(ctx context.Context) ([]*model.MembershipLevel, error)
	GetMembershipLevelByID(ctx context.Context, id int64) (*model.MembershipLevel, error)
	GetMembershipLevelByLevel(ctx context.Context, level int) (*model.MembershipLevel, error)
	CreateMembershipLevel(ctx context.Context, level *model.MembershipLevel) (int64, error)
	UpdateMembershipLevel(ctx context.Context, level *model.MembershipLevel) error
	DeleteMembershipLevel(ctx context.Context, id int64) error
	
	// 会员权益相关
	GetMembershipBenefits(ctx context.Context, levelID int64) ([]*model.MembershipBenefit, error)
	
	// 升级条件相关
	GetMembershipRequirements(ctx context.Context, levelID int64) ([]*model.MembershipRequirement, error)
	
	// 用户会员相关
	GetUserMembership(ctx context.Context, userID int64) (*model.UserMembership, error)
	CreateUserMembership(ctx context.Context, membership *model.UserMembership) (int64, error)
	UpdateUserMembership(ctx context.Context, membership *model.UserMembership) error
	ExtendUserMembership(ctx context.Context, userID int64, durationDays int) error
	
	// 会员交易相关
	CreateMembershipTransaction(ctx context.Context, transaction *model.MembershipTransaction) (int64, error)
	GetMembershipTransactions(ctx context.Context, userID int64, limit, offset int) ([]*model.MembershipTransaction, int, error)
	UpdateMembershipTransaction(ctx context.Context, orderID string, status string) error
	GetMembershipTransactionByOrderID(ctx context.Context, orderID string) (*model.MembershipTransaction, error)
}

// membershipRepo 会员等级存储库实现
type membershipRepo struct {
	db *db.DB
}

// NewMembershipRepo 创建会员等级存储库
func NewMembershipRepo(db *db.DB) MembershipRepo {
	return &membershipRepo{
		db: db,
	}
}

// GetMembershipLevels 获取所有会员等级
func (r *membershipRepo) GetMembershipLevels(ctx context.Context) ([]*model.MembershipLevel, error) {
	return r.db.GetMembershipLevels(ctx)
}

// GetMembershipLevelByID 根据ID获取会员等级
func (r *membershipRepo) GetMembershipLevelByID(ctx context.Context, id int64) (*model.MembershipLevel, error) {
	return r.db.GetMembershipLevelByID(ctx, id)
}

// GetMembershipLevelByLevel 根据等级值获取会员等级
func (r *membershipRepo) GetMembershipLevelByLevel(ctx context.Context, level int) (*model.MembershipLevel, error) {
	return r.db.GetMembershipLevelByLevel(ctx, level)
}

// CreateMembershipLevel 创建会员等级
func (r *membershipRepo) CreateMembershipLevel(ctx context.Context, level *model.MembershipLevel) (int64, error) {
	return r.db.CreateMembershipLevel(ctx, level)
}

// UpdateMembershipLevel 更新会员等级
func (r *membershipRepo) UpdateMembershipLevel(ctx context.Context, level *model.MembershipLevel) error {
	return r.db.UpdateMembershipLevel(ctx, level)
}

// DeleteMembershipLevel 删除会员等级
func (r *membershipRepo) DeleteMembershipLevel(ctx context.Context, id int64) error {
	return r.db.DeleteMembershipLevel(ctx, id)
}

// GetMembershipBenefits 获取会员等级的权益列表
func (r *membershipRepo) GetMembershipBenefits(ctx context.Context, levelID int64) ([]*model.MembershipBenefit, error) {
	return r.db.GetMembershipBenefits(ctx, levelID)
}

// GetMembershipRequirements 获取会员等级的升级条件
func (r *membershipRepo) GetMembershipRequirements(ctx context.Context, levelID int64) ([]*model.MembershipRequirement, error) {
	return r.db.GetMembershipRequirements(ctx, levelID)
}

// GetUserMembership 获取用户的会员信息
func (r *membershipRepo) GetUserMembership(ctx context.Context, userID int64) (*model.UserMembership, error) {
	return r.db.GetUserMembership(ctx, userID)
}

// CreateUserMembership 创建用户会员记录
func (r *membershipRepo) CreateUserMembership(ctx context.Context, membership *model.UserMembership) (int64, error) {
	return r.db.CreateUserMembership(ctx, membership)
}

// UpdateUserMembership 更新用户会员记录
func (r *membershipRepo) UpdateUserMembership(ctx context.Context, membership *model.UserMembership) error {
	return r.db.UpdateUserMembership(ctx, membership)
}

// ExtendUserMembership 延长用户会员期限
func (r *membershipRepo) ExtendUserMembership(ctx context.Context, userID int64, durationDays int) error {
	return r.db.ExtendUserMembership(ctx, userID, durationDays)
}

// CreateMembershipTransaction 创建会员交易记录
func (r *membershipRepo) CreateMembershipTransaction(ctx context.Context, transaction *model.MembershipTransaction) (int64, error) {
	return r.db.CreateMembershipTransaction(ctx, transaction)
}

// GetMembershipTransactions 获取用户会员交易记录
func (r *membershipRepo) GetMembershipTransactions(ctx context.Context, userID int64, limit, offset int) ([]*model.MembershipTransaction, int, error) {
	return r.db.GetMembershipTransactions(ctx, userID, limit, offset)
}

// UpdateMembershipTransaction 更新会员交易记录状态
func (r *membershipRepo) UpdateMembershipTransaction(ctx context.Context, orderID string, status string) error {
	return r.db.UpdateMembershipTransaction(ctx, orderID, status)
}

// GetMembershipTransactionByOrderID 根据订单号获取交易记录
func (r *membershipRepo) GetMembershipTransactionByOrderID(ctx context.Context, orderID string) (*model.MembershipTransaction, error) {
	// 需要在DB层添加实现
	return nil, nil
} 