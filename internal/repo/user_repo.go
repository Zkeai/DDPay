package repo

import (
	"context"
	"database/sql"

	"github.com/Zkeai/DDPay/internal/model"
)

// UserRepo 用户仓库接口
type UserRepo interface {
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int, error)
}

// userRepo 用户仓库实现
type userRepo struct {
	db *sql.DB
}

// NewUserRepo 创建用户仓库
func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

// GetUserByID 根据ID获取用户
func (r *userRepo) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	query := `SELECT * FROM users WHERE id = ?`
	
	var user model.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.Username,
		&user.Avatar, &user.Role, &user.Status, &user.EmailVerified,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt, &user.LastLoginIP,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT * FROM users WHERE email = ?`
	
	var user model.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Username,
		&user.Avatar, &user.Role, &user.Status, &user.EmailVerified,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt, &user.LastLoginIP,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &user, nil
}

// CreateUser 创建用户
func (r *userRepo) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	query := `
		INSERT INTO users (
			email, password, username, avatar, role, status, email_verified,
			created_at, updated_at, last_login_at, last_login_ip
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	result, err := r.db.ExecContext(
		ctx, query,
		user.Email, user.Password, user.Username, user.Avatar,
		user.Role, user.Status, user.EmailVerified, user.CreatedAt,
		user.UpdatedAt, user.LastLoginAt, user.LastLoginIP,
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

// UpdateUser 更新用户
func (r *userRepo) UpdateUser(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users SET
			email = ?, password = ?, username = ?, avatar = ?,
			role = ?, status = ?, level = ?, email_verified = ?, updated_at = ?,
			last_login_at = ?, last_login_ip = ?
		WHERE id = ?
	`
	
	_, err := r.db.ExecContext(
		ctx, query,
		user.Email, user.Password, user.Username, user.Avatar,
		user.Role, user.Status, user.Level, user.EmailVerified, user.UpdatedAt,
		user.LastLoginAt, user.LastLoginIP, user.ID,
	)
	return err
}

// DeleteUser 删除用户
func (r *userRepo) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ListUsers 获取用户列表
func (r *userRepo) ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int, error) {
	// 获取总数
	countQuery := `SELECT COUNT(*) FROM users`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
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
	
	query := `
		SELECT * FROM users
		ORDER BY id DESC
		LIMIT ?, ?
	`
	
	rows, err := r.db.QueryContext(ctx, query, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var users []*model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.Password, &user.Username,
			&user.Avatar, &user.Role, &user.Status, &user.EmailVerified,
			&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt, &user.LastLoginIP,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}
	
	return users, total, nil
} 