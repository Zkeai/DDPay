package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Zkeai/DDPay/internal/model"
)

// GetSubsiteBalance 获取分站余额
func (r *subsiteRepo) GetSubsiteBalance(ctx context.Context, ownerID int64) (*model.SubsiteBalance, error) {
	query := `SELECT * FROM subsite_balances WHERE owner_id = ?`
	
	var balance model.SubsiteBalance
	err := r.db.QueryRowContext(ctx, query, ownerID).Scan(
		&balance.ID, &balance.OwnerID, &balance.Amount, 
		&balance.CreatedAt, &balance.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// 如果不存在，创建一个新的余额记录
			now := time.Now()
			newBalance := &model.SubsiteBalance{
				OwnerID:   ownerID,
				Amount:    0,
				CreatedAt: now,
				UpdatedAt: now,
			}
			
			insertQuery := `
				INSERT INTO subsite_balances (owner_id, amount, created_at, updated_at)
				VALUES (?, ?, ?, ?)
			`
			result, err := r.db.ExecContext(
				ctx, insertQuery,
				newBalance.OwnerID, newBalance.Amount, newBalance.CreatedAt, newBalance.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}
			
			id, err := result.LastInsertId()
			if err != nil {
				return nil, err
			}
			
			newBalance.ID = id
			return newBalance, nil
		}
		return nil, err
	}

	return &balance, nil
}

// UpdateSubsiteBalance 更新分站余额
func (r *subsiteRepo) UpdateSubsiteBalance(ctx context.Context, balance *model.SubsiteBalance) error {
	query := `
		UPDATE subsite_balances SET
			amount = ?, updated_at = ?
		WHERE id = ?
	`
	balance.UpdatedAt = time.Now()
	
	_, err := r.db.ExecContext(
		ctx, query,
		balance.Amount, balance.UpdatedAt, balance.ID,
	)
	return err
}

// CreateSubsiteBalanceLog 创建分站余额变动记录
func (r *subsiteRepo) CreateSubsiteBalanceLog(ctx context.Context, log *model.SubsiteBalanceLog) error {
	query := `
		INSERT INTO subsite_balance_logs (
			owner_id, order_id, amount, before_balance, after_balance,
			type, remark, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	log.CreatedAt = now

	// 处理可能为空的订单ID
	var orderID sql.NullInt64
	if log.OrderID > 0 {
		orderID.Int64 = log.OrderID
		orderID.Valid = true
	}

	_, err := r.db.ExecContext(
		ctx, query,
		log.OwnerID, orderID, log.Amount, log.BeforeBalance,
		log.AfterBalance, log.Type, log.Remark, log.CreatedAt,
	)
	return err
}

// ListSubsiteBalanceLogs 获取分站余额变动记录列表
func (r *subsiteRepo) ListSubsiteBalanceLogs(ctx context.Context, ownerID int64, page, pageSize int) ([]*model.SubsiteBalanceLog, int, error) {
	whereClause := "WHERE owner_id = ?"
	args := []interface{}{ownerID}
	
	// 获取总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM subsite_balance_logs %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	
	query := fmt.Sprintf(`
		SELECT * FROM subsite_balance_logs %s
		ORDER BY id DESC
		LIMIT ?, ?
	`, whereClause)
	args = append(args, offset, pageSize)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var logs []*model.SubsiteBalanceLog
	for rows.Next() {
		var log model.SubsiteBalanceLog
		var orderID sql.NullInt64
		
		err := rows.Scan(
			&log.ID, &log.OwnerID, &orderID, &log.Amount,
			&log.BeforeBalance, &log.AfterBalance, &log.Type,
			&log.Remark, &log.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		
		// 处理可能为空的字段
		if orderID.Valid {
			log.OrderID = orderID.Int64
		}
		
		logs = append(logs, &log)
	}
	
	return logs, total, nil
}

// CreateSubsiteWithdrawal 创建分站提现申请
func (r *subsiteRepo) CreateSubsiteWithdrawal(ctx context.Context, withdrawal *model.SubsiteWithdrawal) (int64, error) {
	query := `
		INSERT INTO subsite_withdrawals (
			owner_id, amount, status, account_type, account_name,
			account_no, remark, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	withdrawal.CreatedAt = now
	withdrawal.UpdatedAt = now

	result, err := r.db.ExecContext(
		ctx, query,
		withdrawal.OwnerID, withdrawal.Amount, withdrawal.Status,
		withdrawal.AccountType, withdrawal.AccountName, withdrawal.AccountNo,
		withdrawal.Remark, withdrawal.CreatedAt, withdrawal.UpdatedAt,
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

// GetSubsiteWithdrawalByID 根据ID获取分站提现申请
func (r *subsiteRepo) GetSubsiteWithdrawalByID(ctx context.Context, id int64) (*model.SubsiteWithdrawal, error) {
	query := `SELECT * FROM subsite_withdrawals WHERE id = ?`
	
	var withdrawal model.SubsiteWithdrawal
	var processedAt sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&withdrawal.ID, &withdrawal.OwnerID, &withdrawal.Amount,
		&withdrawal.Status, &withdrawal.AccountType, &withdrawal.AccountName,
		&withdrawal.AccountNo, &withdrawal.Remark, &withdrawal.AdminRemark,
		&processedAt, &withdrawal.CreatedAt, &withdrawal.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// 处理可能为空的时间字段
	if processedAt.Valid {
		withdrawal.ProcessedAt = processedAt.Time
	}

	return &withdrawal, nil
}

// UpdateSubsiteWithdrawal 更新分站提现申请
func (r *subsiteRepo) UpdateSubsiteWithdrawal(ctx context.Context, withdrawal *model.SubsiteWithdrawal) error {
	query := `
		UPDATE subsite_withdrawals SET
			status = ?, admin_remark = ?, processed_at = ?, updated_at = ?
		WHERE id = ?
	`
	withdrawal.UpdatedAt = time.Now()
	
	// 处理可能为空的时间字段
	var processedAt sql.NullTime
	if !withdrawal.ProcessedAt.IsZero() {
		processedAt.Time = withdrawal.ProcessedAt
		processedAt.Valid = true
	}
	
	_, err := r.db.ExecContext(
		ctx, query,
		withdrawal.Status, withdrawal.AdminRemark, 
		processedAt, withdrawal.UpdatedAt, withdrawal.ID,
	)
	return err
}

// ListSubsiteWithdrawals 获取分站提现申请列表
func (r *subsiteRepo) ListSubsiteWithdrawals(ctx context.Context, ownerID int64, page, pageSize int, status int) ([]*model.SubsiteWithdrawal, int, error) {
	whereClause := "WHERE owner_id = ?"
	args := []interface{}{ownerID}
	
	if status != -1 {
		whereClause += " AND status = ?"
		args = append(args, status)
	}
	
	// 获取总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM subsite_withdrawals %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	
	query := fmt.Sprintf(`
		SELECT * FROM subsite_withdrawals %s
		ORDER BY id DESC
		LIMIT ?, ?
	`, whereClause)
	args = append(args, offset, pageSize)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var withdrawals []*model.SubsiteWithdrawal
	for rows.Next() {
		var withdrawal model.SubsiteWithdrawal
		var processedAt sql.NullTime
		
		err := rows.Scan(
			&withdrawal.ID, &withdrawal.OwnerID, &withdrawal.Amount,
			&withdrawal.Status, &withdrawal.AccountType, &withdrawal.AccountName,
			&withdrawal.AccountNo, &withdrawal.Remark, &withdrawal.AdminRemark,
			&processedAt, &withdrawal.CreatedAt, &withdrawal.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		
		// 处理可能为空的时间字段
		if processedAt.Valid {
			withdrawal.ProcessedAt = processedAt.Time
		}
		
		withdrawals = append(withdrawals, &withdrawal)
	}
	
	return withdrawals, total, nil
} 