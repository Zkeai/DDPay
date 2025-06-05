package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Zkeai/DDPay/internal/model"
	"github.com/Zkeai/DDPay/internal/repo"
)

// MembershipService 会员服务接口
type MembershipService interface {
	// 会员等级相关
	GetMembershipLevels(ctx context.Context) ([]*model.MembershipLevel, error)
	GetMembershipLevelByID(ctx context.Context, id int64) (*model.MembershipLevel, error)
	GetMembershipLevelByLevel(ctx context.Context, level int) (*model.MembershipLevel, error)
	CreateMembershipLevel(ctx context.Context, level *model.MembershipLevel) (int64, error)
	UpdateMembershipLevel(ctx context.Context, level *model.MembershipLevel) error
	DeleteMembershipLevel(ctx context.Context, id int64) error
	
	// 用户会员相关
	GetUserMembership(ctx context.Context, userID int64) (*model.UserMembership, error)
	GetUserMembershipWithLevel(ctx context.Context, userID int64) (*model.UserMembership, *model.MembershipLevel, error)
	PurchaseMembership(ctx context.Context, userID, levelID int64, paymentMethod string, amount float64) (string, error)
	UpgradeMembership(ctx context.Context, userID, levelID int64, paymentMethod string, amount float64) (string, error)
	RenewMembership(ctx context.Context, userID int64, durationDays int, paymentMethod string, amount float64) (string, error)
	
	// 交易记录相关
	GetMembershipTransactions(ctx context.Context, userID int64, page, pageSize int) ([]*model.MembershipTransaction, int, error)
	UpdateTransactionStatus(ctx context.Context, orderID, status string) error
	
	// 检查用户是否满足会员升级条件
	CheckMembershipRequirements(ctx context.Context, userID, levelID int64) (bool, map[string]interface{}, error)
}

// membershipService 会员服务实现
type membershipService struct {
	svc *Service
	repo *repo.Repo
}

// NewMembershipService 创建会员服务
func NewMembershipService(svc *Service) MembershipService {
	return &membershipService{
		svc: svc,
		repo: svc.repo,
	}
}

// GetMembershipLevels 获取所有会员等级
func (s *membershipService) GetMembershipLevels(ctx context.Context) ([]*model.MembershipLevel, error) {
	return s.repo.GetMembershipLevels(ctx)
}

// GetMembershipLevelByID 根据ID获取会员等级
func (s *membershipService) GetMembershipLevelByID(ctx context.Context, id int64) (*model.MembershipLevel, error) {
	return s.repo.GetMembershipLevelByID(ctx, id)
}

// GetMembershipLevelByLevel 根据等级值获取会员等级
func (s *membershipService) GetMembershipLevelByLevel(ctx context.Context, level int) (*model.MembershipLevel, error) {
	return s.repo.GetMembershipLevelByLevel(ctx, level)
}

// GetUserMembership 获取用户会员信息
func (s *membershipService) GetUserMembership(ctx context.Context, userID int64) (*model.UserMembership, error) {
	return s.repo.GetUserMembership(ctx, userID)
}

// GetUserMembershipWithLevel 获取用户会员信息及对应等级详情
func (s *membershipService) GetUserMembershipWithLevel(ctx context.Context, userID int64) (*model.UserMembership, *model.MembershipLevel, error) {
	membership, err := s.repo.GetUserMembership(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	
	// 如果用户没有会员信息，则返回默认的免费会员等级（1级）
	if membership == nil {
		defaultLevel, err := s.repo.GetMembershipLevelByLevel(ctx, 1)
		if err != nil {
			return nil, nil, err
		}
		return nil, defaultLevel, nil
	}
	
	// 获取对应的会员等级信息
	level, err := s.repo.GetMembershipLevelByID(ctx, membership.LevelID)
	if err != nil {
		return membership, nil, err
	}
	
	return membership, level, nil
}

// 生成订单号
func (s *membershipService) generateOrderID(userID int64, transactionType string) string {
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("MEM%s%d%s", transactionType[:3], userID, timestamp)
}

// PurchaseMembership 购买会员
func (s *membershipService) PurchaseMembership(ctx context.Context, userID, levelID int64, paymentMethod string, amount float64) (string, error) {
	// 检查用户是否存在
	user, err := s.svc.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("用户不存在")
	}
	
	// 检查会员等级是否存在
	_, err = s.repo.GetMembershipLevelByID(ctx, levelID)
	if err != nil {
		return "", err
	}
	
	// 检查用户当前会员状态
	currentMembership, err := s.repo.GetUserMembership(ctx, userID)
	if err != nil {
		return "", err
	}
	
	// 如果用户已经是该等级会员，则返回错误
	if currentMembership != nil && currentMembership.LevelID == levelID {
		return "", fmt.Errorf("您已经是该等级会员")
	}
	
	// 生成订单号
	orderID := s.generateOrderID(userID, "PURCHASE")
	
	// 创建交易记录
	transaction := &model.MembershipTransaction{
		UserID:          userID,
		LevelID:         levelID,
		Amount:          amount,
		TransactionType: "purchase",
		PaymentMethod:   paymentMethod,
		Status:          "pending", // 初始状态为待支付
		OrderID:         orderID,
	}
	
	_, err = s.repo.CreateMembershipTransaction(ctx, transaction)
	if err != nil {
		return "", err
	}
	
	// 如果金额为0（免费会员），直接激活会员
	if amount == 0 {
		// 创建会员记录
		endDate := time.Now().AddDate(100, 0, 0) // 设置很长的有效期，相当于永久
		membership := &model.UserMembership{
			UserID:         userID,
			LevelID:        levelID,
			StartDate:      time.Now(),
			EndDate:        &endDate,
			IsActive:       true,
			PurchaseAmount: amount,
		}
		
		_, err = s.repo.CreateUserMembership(ctx, membership)
		if err != nil {
			return "", err
		}
		
		// 更新交易状态为已完成
		err = s.repo.UpdateMembershipTransaction(ctx, orderID, "completed")
		if err != nil {
			return "", err
		}
	}
	
	return orderID, nil
}

// UpgradeMembership 升级会员
func (s *membershipService) UpgradeMembership(ctx context.Context, userID, levelID int64, paymentMethod string, amount float64) (string, error) {
	// 检查用户是否存在
	user, err := s.svc.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("用户不存在")
	}
	
	// 检查会员等级是否存在
	newLevel, err := s.repo.GetMembershipLevelByID(ctx, levelID)
	if err != nil {
		return "", err
	}
	
	// 检查用户当前会员状态
	currentMembership, err := s.repo.GetUserMembership(ctx, userID)
	if err != nil {
		return "", err
	}
	
	// 如果用户没有会员或会员已过期，则调用购买会员方法
	if currentMembership == nil {
		return s.PurchaseMembership(ctx, userID, levelID, paymentMethod, amount)
	}
	
	// 获取当前会员等级
	currentLevel, err := s.repo.GetMembershipLevelByID(ctx, currentMembership.LevelID)
	if err != nil {
		return "", err
	}
	
	// 如果新等级小于等于当前等级，则返回错误
	if newLevel.Level <= currentLevel.Level {
		return "", fmt.Errorf("不能降级或升级到相同等级")
	}
	
	// 生成订单号
	orderID := s.generateOrderID(userID, "UPGRADE")
	
	// 创建交易记录
	transaction := &model.MembershipTransaction{
		UserID:          userID,
		LevelID:         levelID,
		Amount:          amount,
		TransactionType: "upgrade",
		PaymentMethod:   paymentMethod,
		Status:          "pending", // 初始状态为待支付
		OrderID:         orderID,
	}
	
	_, err = s.repo.CreateMembershipTransaction(ctx, transaction)
	if err != nil {
		return "", err
	}
	
	return orderID, nil
}

// RenewMembership 续费会员
func (s *membershipService) RenewMembership(ctx context.Context, userID int64, durationDays int, paymentMethod string, amount float64) (string, error) {
	// 检查用户是否存在
	user, err := s.svc.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("用户不存在")
	}
	
	// 检查用户当前会员状态
	currentMembership, err := s.repo.GetUserMembership(ctx, userID)
	if err != nil {
		return "", err
	}
	
	// 如果用户没有会员，则返回错误
	if currentMembership == nil {
		return "", fmt.Errorf("您还不是会员，请先购买会员")
	}
	
	// 生成订单号
	orderID := s.generateOrderID(userID, "RENEW")
	
	// 创建交易记录
	transaction := &model.MembershipTransaction{
		UserID:          userID,
		LevelID:         currentMembership.LevelID,
		Amount:          amount,
		TransactionType: "renew",
		PaymentMethod:   paymentMethod,
		Status:          "pending", // 初始状态为待支付
		OrderID:         orderID,
	}
	
	_, err = s.repo.CreateMembershipTransaction(ctx, transaction)
	if err != nil {
		return "", err
	}
	
	return orderID, nil
}

// GetMembershipTransactions 获取会员交易记录
func (s *membershipService) GetMembershipTransactions(ctx context.Context, userID int64, page, pageSize int) ([]*model.MembershipTransaction, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	
	offset := (page - 1) * pageSize
	return s.repo.GetMembershipTransactions(ctx, userID, pageSize, offset)
}

// UpdateTransactionStatus 更新交易状态
func (s *membershipService) UpdateTransactionStatus(ctx context.Context, orderID, status string) error {
	// 更新交易状态
	err := s.repo.UpdateMembershipTransaction(ctx, orderID, status)
	if err != nil {
		return err
	}
	
	// 如果交易状态为已完成，则激活/更新会员
	if status == "completed" {
		// 查询交易记录
		// 这里需要额外实现一个根据订单号查询交易记录的方法
		// 简化处理，可以在数据库查询获取交易记录信息
		// 然后根据交易类型(purchase/upgrade/renew)执行相应的会员激活/更新操作
		
		// 此处省略具体实现，实际项目中需要完善
	}
	
	return nil
}

// CheckMembershipRequirements 检查用户是否满足会员升级条件
func (s *membershipService) CheckMembershipRequirements(ctx context.Context, userID, levelID int64) (bool, map[string]interface{}, error) {
	// 获取会员等级
	_, err := s.repo.GetMembershipLevelByID(ctx, levelID)
	if err != nil {
		return false, nil, err
	}
	
	// 获取该等级的升级条件
	requirements, err := s.repo.GetMembershipRequirements(ctx, levelID)
	if err != nil {
		return false, nil, err
	}
	
	// 检查用户是否满足所有条件
	results := make(map[string]interface{})
	allMet := true
	
	// 获取用户信息和统计数据
	// 这里需要实现获取用户相关统计数据的方法，如总订单数、总充值金额、交易额、邀请人数等
	// 简化处理，此处仅做演示
	
	for _, req := range requirements {
		var met bool
		var actual float64
		
		switch req.Type {
		case "register":
			// 注册条件始终满足
			met = true
			actual = 1
		case "payment":
			// 一次性支付金额条件，这个由用户选择支付，始终视为未满足
			met = false
			actual = 0
		case "total_order":
			// 获取用户总订单数
			// 实际项目中需要从数据库查询
			totalOrders := 0 // 示例值
			met = float64(totalOrders) >= req.Value
			actual = float64(totalOrders)
		case "total_payment":
			// 获取用户总充值金额
			// 实际项目中需要从数据库查询
			totalPayment := 0.0 // 示例值
			met = totalPayment >= req.Value
			actual = totalPayment
		case "total_transaction":
			// 获取用户总交易额
			// 实际项目中需要从数据库查询
			totalTransaction := 0.0 // 示例值
			met = totalTransaction >= req.Value
			actual = totalTransaction
		case "invitation":
			// 获取用户邀请人数
			// 实际项目中需要从数据库查询
			invitationCount := 0 // 示例值
			met = float64(invitationCount) >= req.Value
			actual = float64(invitationCount)
		}
		
		results[req.Type] = map[string]interface{}{
			"required": req.Value,
			"actual":   actual,
			"met":      met,
			"description": req.Description,
		}
		
		// 如果有一个条件不满足，则整体不满足
		if !met && req.Type != "payment" { // payment条件特殊处理
			allMet = false
		}
	}
	
	return allMet, results, nil
}

// CreateMembershipLevel 创建会员等级
func (s *membershipService) CreateMembershipLevel(ctx context.Context, level *model.MembershipLevel) (int64, error) {
	// 设置创建和更新时间
	now := time.Now()
	level.CreatedAt = now
	level.UpdatedAt = now
	
	membershipRepo := s.repo.GetMembershipRepo()
	return membershipRepo.CreateMembershipLevel(ctx, level)
}

// UpdateMembershipLevel 更新会员等级
func (s *membershipService) UpdateMembershipLevel(ctx context.Context, level *model.MembershipLevel) error {
	membershipRepo := s.repo.GetMembershipRepo()
	
	// 检查会员等级是否存在
	_, err := membershipRepo.GetMembershipLevelByID(ctx, level.ID)
	if err != nil {
		return fmt.Errorf("会员等级不存在: %w", err)
	}
	
	// 更新时间
	level.UpdatedAt = time.Now()
	
	return membershipRepo.UpdateMembershipLevel(ctx, level)
}

// DeleteMembershipLevel 删除会员等级
func (s *membershipService) DeleteMembershipLevel(ctx context.Context, id int64) error {
	membershipRepo := s.repo.GetMembershipRepo()
	
	// 检查会员等级是否存在
	_, err := membershipRepo.GetMembershipLevelByID(ctx, id)
	if err != nil {
		return fmt.Errorf("会员等级不存在: %w", err)
	}
	
	// 检查是否有用户正在使用该会员等级
	// 实际实现中应该检查用户表中是否有关联该等级的记录
	// 简化处理，这里可以先不实现
	
	return membershipRepo.DeleteMembershipLevel(ctx, id)
} 