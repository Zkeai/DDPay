package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Zkeai/DDPay/internal/model"
)

// SubsiteRepo 分站仓库接口
type SubsiteRepo interface {
	// 分站基本操作
	CreateSubsite(ctx context.Context, subsite *model.Subsite) (int64, error)
	GetSubsiteByID(ctx context.Context, id int64) (*model.Subsite, error)
	GetSubsiteByOwnerID(ctx context.Context, ownerID int64) (*model.Subsite, error)
	GetSubsiteByDomain(ctx context.Context, domain string) (*model.Subsite, error)
	GetSubsiteBySubdomain(ctx context.Context, subdomain string) (*model.Subsite, error)
	UpdateSubsite(ctx context.Context, subsite *model.Subsite) error
	DeleteSubsite(ctx context.Context, id int64) error
	ListSubsites(ctx context.Context, page, pageSize int, status int) ([]*model.Subsite, int, error)
	ListSubsitesByOwnerID(ctx context.Context, ownerID int64) ([]*model.Subsite, int, error)

	// 分站JSON配置操作
	SaveSubsiteConfig(ctx context.Context, config *model.SubsiteConfig) error
	GetSubsiteConfig(ctx context.Context, subsiteID int64) (*model.SubsiteConfig, error)
	
	// 分站商品操作
	CreateSubsiteProduct(ctx context.Context, product *model.SubsiteProduct) (int64, error)
	GetSubsiteProductByID(ctx context.Context, id int64) (*model.SubsiteProduct, error)
	UpdateSubsiteProduct(ctx context.Context, product *model.SubsiteProduct) error
	ListSubsiteProducts(ctx context.Context, subsiteID int64, page, pageSize int, status int) ([]*model.SubsiteProduct, int, error)
	
	// 分站订单操作
	CreateSubsiteOrder(ctx context.Context, order *model.SubsiteOrder) (int64, error)
	GetSubsiteOrderByID(ctx context.Context, id int64) (*model.SubsiteOrder, error)
	GetSubsiteOrderByOrderNo(ctx context.Context, orderNo string) (*model.SubsiteOrder, error)
	UpdateSubsiteOrder(ctx context.Context, order *model.SubsiteOrder) error
	ListSubsiteOrders(ctx context.Context, subsiteID int64, page, pageSize int, status int) ([]*model.SubsiteOrder, int, error)
	
	// 分站余额操作
	GetSubsiteBalance(ctx context.Context, ownerID int64) (*model.SubsiteBalance, error)
	UpdateSubsiteBalance(ctx context.Context, balance *model.SubsiteBalance) error
	CreateSubsiteBalanceLog(ctx context.Context, log *model.SubsiteBalanceLog) error
	ListSubsiteBalanceLogs(ctx context.Context, ownerID int64, page, pageSize int) ([]*model.SubsiteBalanceLog, int, error)
	
	// 分站提现操作
	CreateSubsiteWithdrawal(ctx context.Context, withdrawal *model.SubsiteWithdrawal) (int64, error)
	GetSubsiteWithdrawalByID(ctx context.Context, id int64) (*model.SubsiteWithdrawal, error)
	UpdateSubsiteWithdrawal(ctx context.Context, withdrawal *model.SubsiteWithdrawal) error
	ListSubsiteWithdrawals(ctx context.Context, ownerID int64, page, pageSize int, status int) ([]*model.SubsiteWithdrawal, int, error)
}

// subsiteRepo 分站仓库实现
type subsiteRepo struct {
	db *sql.DB
}

// NewSubsiteRepo 创建分站仓库
func NewSubsiteRepo(db *sql.DB) SubsiteRepo {
	return &subsiteRepo{
		db: db,
	}
}

// CreateSubsite 创建分站
func (r *subsiteRepo) CreateSubsite(ctx context.Context, subsite *model.Subsite) (int64, error) {
	query := `
		INSERT INTO subsites (
			owner_id, name, domain, subdomain, logo, description, 
			theme, status, commission_rate, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	subsite.CreatedAt = now
	subsite.UpdatedAt = now

	result, err := r.db.ExecContext(
		ctx, query,
		subsite.OwnerID, subsite.Name, subsite.Domain, subsite.Subdomain,
		subsite.Logo, subsite.Description, subsite.Theme, subsite.Status,
		subsite.CommissionRate, subsite.CreatedAt, subsite.UpdatedAt,
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

// GetSubsiteByID 根据ID获取分站
func (r *subsiteRepo) GetSubsiteByID(ctx context.Context, id int64) (*model.Subsite, error) {
	query := `SELECT * FROM subsites WHERE id = ?`
	
	var subsite model.Subsite
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&subsite.ID, &subsite.OwnerID, &subsite.Name, &subsite.Domain,
		&subsite.Subdomain, &subsite.Logo, &subsite.Description, &subsite.Theme,
		&subsite.Status, &subsite.CommissionRate, &subsite.CreatedAt, &subsite.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &subsite, nil
}

// GetSubsiteByOwnerID 根据所有者ID获取分站
func (r *subsiteRepo) GetSubsiteByOwnerID(ctx context.Context, ownerID int64) (*model.Subsite, error) {
	query := `SELECT * FROM subsites WHERE owner_id = ? LIMIT 1`
	
	var subsite model.Subsite
	err := r.db.QueryRowContext(ctx, query, ownerID).Scan(
		&subsite.ID, &subsite.OwnerID, &subsite.Name, &subsite.Domain,
		&subsite.Subdomain, &subsite.Logo, &subsite.Description, &subsite.Theme,
		&subsite.Status, &subsite.CommissionRate, &subsite.CreatedAt, &subsite.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &subsite, nil
}

// GetSubsiteByDomain 根据域名获取分站
func (r *subsiteRepo) GetSubsiteByDomain(ctx context.Context, domain string) (*model.Subsite, error) {
	query := `SELECT * FROM subsites WHERE domain = ? LIMIT 1`
	
	var subsite model.Subsite
	err := r.db.QueryRowContext(ctx, query, domain).Scan(
		&subsite.ID, &subsite.OwnerID, &subsite.Name, &subsite.Domain,
		&subsite.Subdomain, &subsite.Logo, &subsite.Description, &subsite.Theme,
		&subsite.Status, &subsite.CommissionRate, &subsite.CreatedAt, &subsite.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &subsite, nil
}

// GetSubsiteBySubdomain 根据子域名获取分站
func (r *subsiteRepo) GetSubsiteBySubdomain(ctx context.Context, subdomain string) (*model.Subsite, error) {
	query := `SELECT * FROM subsites WHERE subdomain = ? LIMIT 1`
	
	var subsite model.Subsite
	err := r.db.QueryRowContext(ctx, query, subdomain).Scan(
		&subsite.ID, &subsite.OwnerID, &subsite.Name, &subsite.Domain,
		&subsite.Subdomain, &subsite.Logo, &subsite.Description, &subsite.Theme,
		&subsite.Status, &subsite.CommissionRate, &subsite.CreatedAt, &subsite.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &subsite, nil
}

// UpdateSubsite 更新分站
func (r *subsiteRepo) UpdateSubsite(ctx context.Context, subsite *model.Subsite) error {
	query := `
		UPDATE subsites SET
			name = ?, domain = ?, subdomain = ?, logo = ?,
			description = ?, theme = ?, status = ?, commission_rate = ?,
			updated_at = ?
		WHERE id = ?
	`
	subsite.UpdatedAt = time.Now()
	
	_, err := r.db.ExecContext(
		ctx, query,
		subsite.Name, subsite.Domain, subsite.Subdomain, subsite.Logo,
		subsite.Description, subsite.Theme, subsite.Status, subsite.CommissionRate,
		subsite.UpdatedAt, subsite.ID,
	)
	return err
}

// ListSubsites 获取分站列表
func (r *subsiteRepo) ListSubsites(ctx context.Context, page, pageSize int, status int) ([]*model.Subsite, int, error) {
	whereClause := ""
	args := []interface{}{}
	
	if status != -1 {
		whereClause = "WHERE status = ?"
		args = append(args, status)
	}
	
	// 获取总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM subsites %s", whereClause)
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
		SELECT * FROM subsites %s
		ORDER BY id DESC
		LIMIT ?, ?
	`, whereClause)
	args = append(args, offset, pageSize)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var subsites []*model.Subsite
	for rows.Next() {
		var subsite model.Subsite
		err := rows.Scan(
			&subsite.ID, &subsite.OwnerID, &subsite.Name, &subsite.Domain,
			&subsite.Subdomain, &subsite.Logo, &subsite.Description, &subsite.Theme,
			&subsite.Status, &subsite.CommissionRate, &subsite.CreatedAt, &subsite.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		subsites = append(subsites, &subsite)
	}
	
	return subsites, total, nil
}

// DeleteSubsite 删除分站
func (r *subsiteRepo) DeleteSubsite(ctx context.Context, id int64) error {
	query := `DELETE FROM subsites WHERE id = ?`
	
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// SaveSubsiteConfig 保存分站JSON配置
func (r *subsiteRepo) SaveSubsiteConfig(ctx context.Context, config *model.SubsiteConfig) error {
	query := `
		INSERT INTO subsite_configs (subsite_id, config, created_at, updated_at)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		config = ?, updated_at = ?
	`
	now := time.Now()
	
	_, err := r.db.ExecContext(
		ctx, query,
		config.SubsiteID, config.Config, now, now,
		config.Config, now,
	)
	return err
}

// GetSubsiteConfig 获取分站JSON配置
func (r *subsiteRepo) GetSubsiteConfig(ctx context.Context, subsiteID int64) (*model.SubsiteConfig, error) {
	query := `SELECT * FROM subsite_configs WHERE subsite_id = ? LIMIT 1`
	
	var config model.SubsiteConfig
	err := r.db.QueryRowContext(ctx, query, subsiteID).Scan(
		&config.ID, &config.SubsiteID, &config.Config, 
		&config.CreatedAt, &config.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &config, nil
}

// ListSubsitesByOwnerID 获取指定所有者的所有分站
func (r *subsiteRepo) ListSubsitesByOwnerID(ctx context.Context, ownerID int64) ([]*model.Subsite, int, error) {
	query := `SELECT * FROM subsites WHERE owner_id = ?`
	countQuery := `SELECT COUNT(*) FROM subsites WHERE owner_id = ?`
	
	// 获取总数
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, ownerID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	
	// 获取分站列表
	rows, err := r.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var subsites []*model.Subsite
	for rows.Next() {
		var subsite model.Subsite
		err := rows.Scan(
			&subsite.ID, &subsite.OwnerID, &subsite.Name, &subsite.Domain,
			&subsite.Subdomain, &subsite.Logo, &subsite.Description, &subsite.Theme,
			&subsite.Status, &subsite.CommissionRate, &subsite.CreatedAt, &subsite.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		subsites = append(subsites, &subsite)
	}
	
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	
	return subsites, total, nil
} 