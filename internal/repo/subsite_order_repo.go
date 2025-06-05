package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Zkeai/DDPay/internal/model"
)

// CreateSubsiteOrder 创建分站订单
func (r *subsiteRepo) CreateSubsiteOrder(ctx context.Context, order *model.SubsiteOrder) (int64, error) {
	query := `
		INSERT INTO subsite_orders (
			order_no, subsite_id, user_id, product_id, quantity, 
			amount, commission, status, pay_time, complete_time,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now

	// 处理可能为空的时间字段和用户ID
	var payTime, completeTime sql.NullTime
	var userID sql.NullInt64
	
	if !order.PayTime.IsZero() {
		payTime.Time = order.PayTime
		payTime.Valid = true
	}
	if !order.CompleteTime.IsZero() {
		completeTime.Time = order.CompleteTime
		completeTime.Valid = true
	}
	if order.UserID > 0 {
		userID.Int64 = order.UserID
		userID.Valid = true
	}

	result, err := r.db.ExecContext(
		ctx, query,
		order.OrderNo, order.SubsiteID, userID, order.ProductID,
		order.Quantity, order.Amount, order.Commission, order.Status,
		payTime, completeTime, order.CreatedAt, order.UpdatedAt,
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

// GetSubsiteOrderByID 根据ID获取分站订单
func (r *subsiteRepo) GetSubsiteOrderByID(ctx context.Context, id int64) (*model.SubsiteOrder, error) {
	query := `SELECT * FROM subsite_orders WHERE id = ?`
	
	var order model.SubsiteOrder
	var payTime, completeTime sql.NullTime
	var userID sql.NullInt64
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID, &order.OrderNo, &order.SubsiteID, &userID, &order.ProductID,
		&order.Quantity, &order.Amount, &order.Commission, &order.Status,
		&payTime, &completeTime, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// 处理可能为空的字段
	if userID.Valid {
		order.UserID = userID.Int64
	}
	if payTime.Valid {
		order.PayTime = payTime.Time
	}
	if completeTime.Valid {
		order.CompleteTime = completeTime.Time
	}

	return &order, nil
}

// GetSubsiteOrderByOrderNo 根据订单号获取分站订单
func (r *subsiteRepo) GetSubsiteOrderByOrderNo(ctx context.Context, orderNo string) (*model.SubsiteOrder, error) {
	query := `SELECT * FROM subsite_orders WHERE order_no = ?`
	
	var order model.SubsiteOrder
	var payTime, completeTime sql.NullTime
	var userID sql.NullInt64
	
	err := r.db.QueryRowContext(ctx, query, orderNo).Scan(
		&order.ID, &order.OrderNo, &order.SubsiteID, &userID, &order.ProductID,
		&order.Quantity, &order.Amount, &order.Commission, &order.Status,
		&payTime, &completeTime, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// 处理可能为空的字段
	if userID.Valid {
		order.UserID = userID.Int64
	}
	if payTime.Valid {
		order.PayTime = payTime.Time
	}
	if completeTime.Valid {
		order.CompleteTime = completeTime.Time
	}

	return &order, nil
}

// UpdateSubsiteOrder 更新分站订单
func (r *subsiteRepo) UpdateSubsiteOrder(ctx context.Context, order *model.SubsiteOrder) error {
	query := `
		UPDATE subsite_orders SET
			quantity = ?, amount = ?, commission = ?, status = ?,
			pay_time = ?, complete_time = ?, updated_at = ?
		WHERE id = ?
	`
	order.UpdatedAt = time.Now()
	
	// 处理可能为空的时间字段
	var payTime, completeTime sql.NullTime
	if !order.PayTime.IsZero() {
		payTime.Time = order.PayTime
		payTime.Valid = true
	}
	if !order.CompleteTime.IsZero() {
		completeTime.Time = order.CompleteTime
		completeTime.Valid = true
	}
	
	_, err := r.db.ExecContext(
		ctx, query,
		order.Quantity, order.Amount, order.Commission, order.Status,
		payTime, completeTime, order.UpdatedAt, order.ID,
	)
	return err
}

// ListSubsiteOrders 获取分站订单列表
func (r *subsiteRepo) ListSubsiteOrders(ctx context.Context, subsiteID int64, page, pageSize int, status int) ([]*model.SubsiteOrder, int, error) {
	whereClause := "WHERE subsite_id = ?"
	args := []interface{}{subsiteID}
	
	if status != -1 {
		whereClause += " AND status = ?"
		args = append(args, status)
	}
	
	// 获取总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM subsite_orders %s", whereClause)
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
		SELECT * FROM subsite_orders %s
		ORDER BY id DESC
		LIMIT ?, ?
	`, whereClause)
	args = append(args, offset, pageSize)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var orders []*model.SubsiteOrder
	for rows.Next() {
		var order model.SubsiteOrder
		var payTime, completeTime sql.NullTime
		var userID sql.NullInt64
		
		err := rows.Scan(
			&order.ID, &order.OrderNo, &order.SubsiteID, &userID, &order.ProductID,
			&order.Quantity, &order.Amount, &order.Commission, &order.Status,
			&payTime, &completeTime, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		
		// 处理可能为空的字段
		if userID.Valid {
			order.UserID = userID.Int64
		}
		if payTime.Valid {
			order.PayTime = payTime.Time
		}
		if completeTime.Valid {
			order.CompleteTime = completeTime.Time
		}
		
		orders = append(orders, &order)
	}
	
	return orders, total, nil
} 