package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Zkeai/DDPay/internal/model"
)

// 会员相关错误
var (
	ErrMembershipLevelNotFound  = errors.New("会员等级不存在")
	ErrMembershipAlreadyExists  = errors.New("用户已经拥有该会员等级")
	ErrInvalidMembershipPeriod  = errors.New("无效的会员期限")
)

// GetMembershipLevels 获取所有会员等级
func (db *DB) GetMembershipLevels(ctx context.Context) ([]*model.MembershipLevel, error) {
	query := `SELECT id, name, level, icon, price, description, discount_rate, 
              max_subsites, custom_service_access, vip_group_access, priority, created_at, updated_at 
              FROM membership_levels ORDER BY level ASC`
	
	rows, err := db.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var levels []*model.MembershipLevel
	for rows.Next() {
		level := &model.MembershipLevel{}
		var customServiceAccess, vipGroupAccess int
		
		err := rows.Scan(
			&level.ID, &level.Name, &level.Level, &level.Icon, &level.Price,
			&level.Description, &level.DiscountRate, &level.MaxSubsites,
			&customServiceAccess, &vipGroupAccess, &level.Priority,
			&level.CreatedAt, &level.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		level.CustomServiceAccess = customServiceAccess == 1
		level.VIPGroupAccess = vipGroupAccess == 1
		
		// 获取会员权益
		benefits, err := db.GetMembershipBenefits(ctx, level.ID)
		if err != nil {
			return nil, err
		}
		level.Benefits = benefits
		
		// 获取升级条件
		requirements, err := db.GetMembershipRequirements(ctx, level.ID)
		if err != nil {
			return nil, err
		}
		level.Requirements = requirements
		
		levels = append(levels, level)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return levels, nil
}

// GetMembershipLevelByID 根据ID获取会员等级
func (db *DB) GetMembershipLevelByID(ctx context.Context, id int64) (*model.MembershipLevel, error) {
	query := `SELECT id, name, level, icon, price, description, discount_rate, 
              max_subsites, custom_service_access, vip_group_access, priority, created_at, updated_at 
              FROM membership_levels WHERE id = ?`
	
	level := &model.MembershipLevel{}
	var customServiceAccess, vipGroupAccess int
	
	err := db.db.QueryRow(ctx, query, id).Scan(
		&level.ID, &level.Name, &level.Level, &level.Icon, &level.Price,
		&level.Description, &level.DiscountRate, &level.MaxSubsites,
		&customServiceAccess, &vipGroupAccess, &level.Priority,
		&level.CreatedAt, &level.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMembershipLevelNotFound
		}
		return nil, err
	}
	
	level.CustomServiceAccess = customServiceAccess == 1
	level.VIPGroupAccess = vipGroupAccess == 1
	
	// 获取会员权益
	benefits, err := db.GetMembershipBenefits(ctx, level.ID)
	if err != nil {
		return nil, err
	}
	level.Benefits = benefits
	
	// 获取升级条件
	requirements, err := db.GetMembershipRequirements(ctx, level.ID)
	if err != nil {
		return nil, err
	}
	level.Requirements = requirements
	
	return level, nil
}

// GetMembershipLevelByLevel 根据等级值获取会员等级
func (db *DB) GetMembershipLevelByLevel(ctx context.Context, levelValue int) (*model.MembershipLevel, error) {
	query := `SELECT id, name, level, icon, price, description, discount_rate, 
              max_subsites, custom_service_access, vip_group_access, priority, created_at, updated_at 
              FROM membership_levels WHERE level = ?`
	
	level := &model.MembershipLevel{}
	var customServiceAccess, vipGroupAccess int
	
	err := db.db.QueryRow(ctx, query, levelValue).Scan(
		&level.ID, &level.Name, &level.Level, &level.Icon, &level.Price,
		&level.Description, &level.DiscountRate, &level.MaxSubsites,
		&customServiceAccess, &vipGroupAccess, &level.Priority,
		&level.CreatedAt, &level.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMembershipLevelNotFound
		}
		return nil, err
	}
	
	level.CustomServiceAccess = customServiceAccess == 1
	level.VIPGroupAccess = vipGroupAccess == 1
	
	// 获取会员权益
	benefits, err := db.GetMembershipBenefits(ctx, level.ID)
	if err != nil {
		return nil, err
	}
	level.Benefits = benefits
	
	// 获取升级条件
	requirements, err := db.GetMembershipRequirements(ctx, level.ID)
	if err != nil {
		return nil, err
	}
	level.Requirements = requirements
	
	return level, nil
}

// GetMembershipBenefits 获取会员等级的权益列表
func (db *DB) GetMembershipBenefits(ctx context.Context, levelID int64) ([]*model.MembershipBenefit, error) {
	query := `SELECT id, level_id, title, description, icon, created_at, updated_at 
              FROM membership_benefits WHERE level_id = ? ORDER BY id ASC`
	
	rows, err := db.db.Query(ctx, query, levelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var benefits []*model.MembershipBenefit
	for rows.Next() {
		benefit := &model.MembershipBenefit{}
		err := rows.Scan(
			&benefit.ID, &benefit.LevelID, &benefit.Title,
			&benefit.Description, &benefit.Icon, &benefit.CreatedAt, &benefit.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		benefits = append(benefits, benefit)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return benefits, nil
}

// GetMembershipRequirements 获取会员等级的升级条件
func (db *DB) GetMembershipRequirements(ctx context.Context, levelID int64) ([]*model.MembershipRequirement, error) {
	query := `SELECT id, level_id, type, value, description, created_at, updated_at 
              FROM membership_requirements WHERE level_id = ? ORDER BY id ASC`
	
	rows, err := db.db.Query(ctx, query, levelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var requirements []*model.MembershipRequirement
	for rows.Next() {
		requirement := &model.MembershipRequirement{}
		err := rows.Scan(
			&requirement.ID, &requirement.LevelID, &requirement.Type,
			&requirement.Value, &requirement.Description, &requirement.CreatedAt, &requirement.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, requirement)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return requirements, nil
}

// GetUserMembership 获取用户的会员信息
func (db *DB) GetUserMembership(ctx context.Context, userID int64) (*model.UserMembership, error) {
	query := `SELECT id, user_id, level_id, start_date, end_date, is_active, purchase_amount, created_at, updated_at 
              FROM user_memberships WHERE user_id = ? AND is_active = 1 
              ORDER BY level_id DESC LIMIT 1`
	
	membership := &model.UserMembership{}
	var isActive int
	
	err := db.db.QueryRow(ctx, query, userID).Scan(
		&membership.ID, &membership.UserID, &membership.LevelID,
		&membership.StartDate, &membership.EndDate, &isActive,
		&membership.PurchaseAmount, &membership.CreatedAt, &membership.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 用户没有会员记录，返回nil
		}
		return nil, err
	}
	
	membership.IsActive = isActive == 1
	
	return membership, nil
}

// CreateUserMembership 创建用户会员记录
func (db *DB) CreateUserMembership(ctx context.Context, membership *model.UserMembership) (int64, error) {
	// 检查是否已存在激活的会员记录
	existingMembership, err := db.GetUserMembership(ctx, membership.UserID)
	if err != nil {
		return 0, err
	}
	
	if existingMembership != nil && existingMembership.LevelID == membership.LevelID {
		return 0, ErrMembershipAlreadyExists
	}
	
	// 如果存在其他会员等级，则将其设为非激活
	if existingMembership != nil {
		deactivateQuery := `UPDATE user_memberships SET is_active = 0, updated_at = ? WHERE id = ?`
		_, err = db.db.Exec(ctx, deactivateQuery, time.Now(), existingMembership.ID)
		if err != nil {
			return 0, err
		}
	}
	
	// 设置默认值
	now := time.Now()
	if membership.CreatedAt.IsZero() {
		membership.CreatedAt = now
	}
	if membership.UpdatedAt.IsZero() {
		membership.UpdatedAt = now
	}
	if membership.StartDate.IsZero() {
		membership.StartDate = now
	}
	
	// 插入新记录
	isActive := 0
	if membership.IsActive {
		isActive = 1
	}
	
	query := `INSERT INTO user_memberships (user_id, level_id, start_date, end_date, is_active, purchase_amount, created_at, updated_at) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := db.db.Exec(ctx, query,
		membership.UserID, membership.LevelID, membership.StartDate, membership.EndDate,
		isActive, membership.PurchaseAmount, membership.CreatedAt, membership.UpdatedAt,
	)
	
	if err != nil {
		return 0, err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	
	return id, nil
}

// UpdateUserMembership 更新用户会员记录
func (db *DB) UpdateUserMembership(ctx context.Context, membership *model.UserMembership) error {
	membership.UpdatedAt = time.Now()
	
	isActive := 0
	if membership.IsActive {
		isActive = 1
	}
	
	query := `UPDATE user_memberships SET level_id = ?, start_date = ?, end_date = ?, 
              is_active = ?, purchase_amount = ?, updated_at = ? WHERE id = ?`
	
	_, err := db.db.Exec(ctx, query,
		membership.LevelID, membership.StartDate, membership.EndDate,
		isActive, membership.PurchaseAmount, membership.UpdatedAt, membership.ID,
	)
	
	return err
}

// ExtendUserMembership 延长用户会员期限
func (db *DB) ExtendUserMembership(ctx context.Context, userID int64, durationDays int) error {
	if durationDays <= 0 {
		return ErrInvalidMembershipPeriod
	}
	
	// 获取当前会员记录
	membership, err := db.GetUserMembership(ctx, userID)
	if err != nil {
		return err
	}
	
	if membership == nil {
		return ErrMembershipLevelNotFound
	}
	
	// 计算新的结束日期
	now := time.Now()
	var newEndDate time.Time
	
	if membership.EndDate == nil || membership.EndDate.Before(now) {
		// 如果已过期或无结束日期，从当前时间开始计算
		newEndDate = now.AddDate(0, 0, durationDays)
	} else {
		// 否则在原结束日期基础上延长
		newEndDate = membership.EndDate.AddDate(0, 0, durationDays)
	}
	
	// 更新结束日期
	query := `UPDATE user_memberships SET end_date = ?, updated_at = ? WHERE id = ?`
	_, err = db.db.Exec(ctx, query, newEndDate, now, membership.ID)
	
	return err
}

// CreateMembershipTransaction 创建会员交易记录
func (db *DB) CreateMembershipTransaction(ctx context.Context, transaction *model.MembershipTransaction) (int64, error) {
	// 设置默认值
	now := time.Now()
	if transaction.CreatedAt.IsZero() {
		transaction.CreatedAt = now
	}
	if transaction.UpdatedAt.IsZero() {
		transaction.UpdatedAt = now
	}
	
	query := `INSERT INTO membership_transactions (user_id, level_id, amount, transaction_type, 
              payment_method, status, order_id, created_at, updated_at) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := db.db.Exec(ctx, query,
		transaction.UserID, transaction.LevelID, transaction.Amount,
		transaction.TransactionType, transaction.PaymentMethod,
		transaction.Status, transaction.OrderID,
		transaction.CreatedAt, transaction.UpdatedAt,
	)
	
	if err != nil {
		return 0, err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	
	return id, nil
}

// GetMembershipTransactions 获取用户会员交易记录
func (db *DB) GetMembershipTransactions(ctx context.Context, userID int64, limit, offset int) ([]*model.MembershipTransaction, int, error) {
	// 获取总数
	countQuery := `SELECT COUNT(*) FROM membership_transactions WHERE user_id = ?`
	var total int
	err := db.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	
	// 查询交易记录
	query := `SELECT id, user_id, level_id, amount, transaction_type, payment_method, 
              status, order_id, created_at, updated_at 
              FROM membership_transactions 
              WHERE user_id = ? 
              ORDER BY created_at DESC 
              LIMIT ? OFFSET ?`
	
	rows, err := db.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var transactions []*model.MembershipTransaction
	for rows.Next() {
		transaction := &model.MembershipTransaction{}
		err := rows.Scan(
			&transaction.ID, &transaction.UserID, &transaction.LevelID,
			&transaction.Amount, &transaction.TransactionType, &transaction.PaymentMethod,
			&transaction.Status, &transaction.OrderID, &transaction.CreatedAt, &transaction.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		transactions = append(transactions, transaction)
	}
	
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	
	return transactions, total, nil
}

// UpdateMembershipTransaction 更新会员交易记录状态
func (db *DB) UpdateMembershipTransaction(ctx context.Context, orderID string, status string) error {
	query := `UPDATE membership_transactions SET status = ?, updated_at = ? WHERE order_id = ?`
	_, err := db.db.Exec(ctx, query, status, time.Now(), orderID)
	return err
}

// CreateMembershipLevel 创建会员等级
func (db *DB) CreateMembershipLevel(ctx context.Context, level *model.MembershipLevel) (int64, error) {
	now := time.Now()
	if level.CreatedAt.IsZero() {
		level.CreatedAt = now
	}
	if level.UpdatedAt.IsZero() {
		level.UpdatedAt = now
	}

	result, err := db.db.Exec(ctx,
		"INSERT INTO membership_levels (name, level, icon, price, description, discount_rate, max_subsites, custom_service_access, vip_group_access, priority, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		level.Name, level.Level, level.Icon, level.Price, level.Description, level.DiscountRate, level.MaxSubsites, level.CustomServiceAccess, level.VIPGroupAccess, level.Priority, level.CreatedAt, level.UpdatedAt)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// UpdateMembershipLevel 更新会员等级
func (db *DB) UpdateMembershipLevel(ctx context.Context, level *model.MembershipLevel) error {
	level.UpdatedAt = time.Now()

	_, err := db.db.Exec(ctx,
		"UPDATE membership_levels SET name = ?, level = ?, icon = ?, price = ?, description = ?, discount_rate = ?, max_subsites = ?, custom_service_access = ?, vip_group_access = ?, priority = ?, updated_at = ? WHERE id = ?",
		level.Name, level.Level, level.Icon, level.Price, level.Description, level.DiscountRate, level.MaxSubsites, level.CustomServiceAccess, level.VIPGroupAccess, level.Priority, level.UpdatedAt, level.ID)
	return err
}

// DeleteMembershipLevel 删除会员等级
func (db *DB) DeleteMembershipLevel(ctx context.Context, id int64) error {
	_, err := db.db.Exec(ctx, "DELETE FROM membership_levels WHERE id = ?", id)
	return err
} 