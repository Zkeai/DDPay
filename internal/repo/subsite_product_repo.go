package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Zkeai/DDPay/internal/model"
)

// CreateSubsiteProduct 创建分站商品
func (r *subsiteRepo) CreateSubsiteProduct(ctx context.Context, product *model.SubsiteProduct) (int64, error) {
	query := `
		INSERT INTO subsite_products (
			subsite_id, main_product_id, name, description, price, 
			original_price, stock, image, status, is_time_limited,
			start_time, end_time, sort_order, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	// 处理可能为空的时间字段
	var startTime, endTime sql.NullTime
	if !product.StartTime.IsZero() {
		startTime.Time = product.StartTime
		startTime.Valid = true
	}
	if !product.EndTime.IsZero() {
		endTime.Time = product.EndTime
		endTime.Valid = true
	}

	result, err := r.db.ExecContext(
		ctx, query,
		product.SubsiteID, product.MainProductID, product.Name, product.Description,
		product.Price, product.OriginalPrice, product.Stock, product.Image,
		product.Status, product.IsTimeLimited, startTime, endTime,
		product.SortOrder, product.CreatedAt, product.UpdatedAt,
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

// GetSubsiteProductByID 根据ID获取分站商品
func (r *subsiteRepo) GetSubsiteProductByID(ctx context.Context, id int64) (*model.SubsiteProduct, error) {
	query := `SELECT * FROM subsite_products WHERE id = ?`
	
	var product model.SubsiteProduct
	var startTime, endTime sql.NullTime
	var mainProductID sql.NullInt64
	var originalPrice sql.NullFloat64
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID, &product.SubsiteID, &mainProductID, &product.Name,
		&product.Description, &product.Price, &originalPrice, &product.Stock,
		&product.Image, &product.Status, &product.IsTimeLimited, 
		&startTime, &endTime, &product.SortOrder,
		&product.CreatedAt, &product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// 处理可能为空的字段
	if mainProductID.Valid {
		product.MainProductID = mainProductID.Int64
	}
	if originalPrice.Valid {
		product.OriginalPrice = originalPrice.Float64
	}
	if startTime.Valid {
		product.StartTime = startTime.Time
	}
	if endTime.Valid {
		product.EndTime = endTime.Time
	}

	return &product, nil
}

// UpdateSubsiteProduct 更新分站商品
func (r *subsiteRepo) UpdateSubsiteProduct(ctx context.Context, product *model.SubsiteProduct) error {
	query := `
		UPDATE subsite_products SET
			name = ?, description = ?, price = ?, original_price = ?,
			stock = ?, image = ?, status = ?, is_time_limited = ?,
			start_time = ?, end_time = ?, sort_order = ?, updated_at = ?
		WHERE id = ?
	`
	product.UpdatedAt = time.Now()
	
	// 处理可能为空的时间字段
	var startTime, endTime sql.NullTime
	if !product.StartTime.IsZero() {
		startTime.Time = product.StartTime
		startTime.Valid = true
	}
	if !product.EndTime.IsZero() {
		endTime.Time = product.EndTime
		endTime.Valid = true
	}
	
	_, err := r.db.ExecContext(
		ctx, query,
		product.Name, product.Description, product.Price, product.OriginalPrice,
		product.Stock, product.Image, product.Status, product.IsTimeLimited,
		startTime, endTime, product.SortOrder, product.UpdatedAt, product.ID,
	)
	return err
}

// ListSubsiteProducts 获取分站商品列表
func (r *subsiteRepo) ListSubsiteProducts(ctx context.Context, subsiteID int64, page, pageSize int, status int) ([]*model.SubsiteProduct, int, error) {
	whereClause := "WHERE subsite_id = ?"
	args := []interface{}{subsiteID}
	
	if status != -1 {
		whereClause += " AND status = ?"
		args = append(args, status)
	}
	
	// 获取总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM subsite_products %s", whereClause)
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
		SELECT * FROM subsite_products %s
		ORDER BY sort_order DESC, id DESC
		LIMIT ?, ?
	`, whereClause)
	args = append(args, offset, pageSize)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var products []*model.SubsiteProduct
	for rows.Next() {
		var product model.SubsiteProduct
		var startTime, endTime sql.NullTime
		var mainProductID sql.NullInt64
		var originalPrice sql.NullFloat64
		
		err := rows.Scan(
			&product.ID, &product.SubsiteID, &mainProductID, &product.Name,
			&product.Description, &product.Price, &originalPrice, &product.Stock,
			&product.Image, &product.Status, &product.IsTimeLimited, 
			&startTime, &endTime, &product.SortOrder,
			&product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		
		// 处理可能为空的字段
		if mainProductID.Valid {
			product.MainProductID = mainProductID.Int64
		}
		if originalPrice.Valid {
			product.OriginalPrice = originalPrice.Float64
		}
		if startTime.Valid {
			product.StartTime = startTime.Time
		}
		if endTime.Valid {
			product.EndTime = endTime.Time
		}
		
		products = append(products, &product)
	}
	
	return products, total, nil
} 