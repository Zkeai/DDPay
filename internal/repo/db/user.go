package db

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/Zkeai/DDPay/internal/model"
)

// 用户相关错误
var (
	ErrUserNotFound      = errors.New("用户不存在")
	ErrEmailAlreadyExists = errors.New("邮箱已被注册")
	ErrInvalidCredentials = errors.New("无效的凭证")
)

// GetUserByID 通过ID获取用户
func (db *DB) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, email, password, username, avatar, role, status, level, email_verified, 
              created_at, updated_at, last_login_at, last_login_ip 
              FROM users WHERE id = ?`
	
	err := db.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.Username, &user.Avatar,
		&user.Role, &user.Status, &user.Level, &user.EmailVerified, &user.CreatedAt,
		&user.UpdatedAt, &user.LastLoginAt, &user.LastLoginIP,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	
	return user, nil
}

// GetUserByEmail 通过邮箱获取用户
func (db *DB) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, email, password, username, avatar, role, status, level, email_verified, 
              created_at, updated_at, last_login_at, last_login_ip 
              FROM users WHERE email = ?`
	
	err := db.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Username, &user.Avatar,
		&user.Role, &user.Status, &user.Level, &user.EmailVerified, &user.CreatedAt,
		&user.UpdatedAt, &user.LastLoginAt, &user.LastLoginIP,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	
	return user, nil
}

// CreateUser 创建用户
func (db *DB) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	// 检查邮箱是否已存在
	exists, err := db.checkEmailExists(ctx, user.Email)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, ErrEmailAlreadyExists
	}
	
	// 设置默认值
	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = now
	}
	if user.Role == "" {
		user.Role = "user" // 默认角色
	}
	
	query := `INSERT INTO users (email, password, username, avatar, role, status, email_verified, 
              created_at, updated_at, last_login_at, last_login_ip) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := db.db.Exec(ctx, query,
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

// UpdateUser 更新用户信息
func (db *DB) UpdateUser(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()
	
	query := `UPDATE users SET username = ?, avatar = ?, status = ?, level = ?, email_verified = ?, 
              updated_at = ? WHERE id = ?`
	
	_, err := db.db.Exec(ctx, query,
		user.Username, user.Avatar, user.Status, user.Level, user.EmailVerified,
		user.UpdatedAt, user.ID,
	)
	
	return err
}

// UpdatePassword 更新用户密码
func (db *DB) UpdatePassword(ctx context.Context, userID int64, password string) error {
	query := `UPDATE users SET password = ?, updated_at = ? WHERE id = ?`
	_, err := db.db.Exec(ctx, query, password, time.Now(), userID)
	return err
}

// UpdateLastLogin 更新最后登录信息
func (db *DB) UpdateLastLogin(ctx context.Context, userID int64, ip string) error {
	query := `UPDATE users SET last_login_at = ?, last_login_ip = ? WHERE id = ?`
	_, err := db.db.Exec(ctx, query, time.Now(), ip, userID)
	return err
}

// checkEmailExists 检查邮箱是否已存在
func (db *DB) checkEmailExists(ctx context.Context, email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = ?`
	err := db.db.QueryRow(ctx, query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateVerificationCode 创建验证码
func (db *DB) CreateVerificationCode(ctx context.Context, code *model.VerificationCode) error {
	query := `INSERT INTO verification_codes (email, code, type, expires_at, used, created_at) 
              VALUES (?, ?, ?, ?, ?, ?)`
	
	_, err := db.db.Exec(ctx, query,
		code.Email, code.Code, code.Type, code.ExpiresAt, code.Used, code.CreatedAt,
	)
	
	return err
}

// GetVerificationCode 获取验证码
func (db *DB) GetVerificationCode(ctx context.Context, email, codeType string) (*model.VerificationCode, error) {
	code := &model.VerificationCode{}
	query := `SELECT id, email, code, type, expires_at, used, created_at 
              FROM verification_codes 
              WHERE email = ? AND type = ? AND used = 0 AND expires_at > ?
              ORDER BY created_at DESC LIMIT 1`
	
	err := db.db.QueryRow(ctx, query, email, codeType, time.Now()).Scan(
		&code.ID, &code.Email, &code.Code, &code.Type,
		&code.ExpiresAt, &code.Used, &code.CreatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("验证码不存在或已过期")
		}
		return nil, err
	}
	
	return code, nil
}

// MarkVerificationCodeAsUsed 标记验证码为已使用
func (db *DB) MarkVerificationCodeAsUsed(ctx context.Context, id int64) error {
	query := `UPDATE verification_codes SET used = 1 WHERE id = ?`
	_, err := db.db.Exec(ctx, query, id)
	return err
}

// CreateOAuthAccount 创建OAuth账号关联
func (db *DB) CreateOAuthAccount(ctx context.Context, account *model.OAuthAccount) error {
	query := `INSERT INTO oauth_accounts (user_id, provider, provider_user_id, provider_username, 
              provider_email, provider_avatar, access_token, refresh_token, token_expires_at, 
              created_at, updated_at) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.db.Exec(ctx, query,
		account.UserID, account.Provider, account.ProviderUserID, account.ProviderUsername,
		account.ProviderEmail, account.ProviderAvatar, account.AccessToken, account.RefreshToken,
		account.TokenExpiresAt, account.CreatedAt, account.UpdatedAt,
	)
	
	return err
}

// GetOAuthAccount 获取OAuth账号关联
func (db *DB) GetOAuthAccount(ctx context.Context, provider, providerUserID string) (*model.OAuthAccount, error) {
	account := &model.OAuthAccount{}
	query := `SELECT id, user_id, provider, provider_user_id, provider_username, provider_email, 
              provider_avatar, access_token, refresh_token, token_expires_at, created_at, updated_at 
              FROM oauth_accounts 
              WHERE provider = ? AND provider_user_id = ?`
	
	err := db.db.QueryRow(ctx, query, provider, providerUserID).Scan(
		&account.ID, &account.UserID, &account.Provider, &account.ProviderUserID,
		&account.ProviderUsername, &account.ProviderEmail, &account.ProviderAvatar,
		&account.AccessToken, &account.RefreshToken, &account.TokenExpiresAt,
		&account.CreatedAt, &account.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 返回nil表示未找到
		}
		return nil, err
	}
	
	return account, nil
}

// UpdateOAuthAccount 更新OAuth账号关联
func (db *DB) UpdateOAuthAccount(ctx context.Context, account *model.OAuthAccount) error {
	query := `UPDATE oauth_accounts SET provider_username = ?, provider_email = ?, provider_avatar = ?, 
              access_token = ?, refresh_token = ?, token_expires_at = ?, updated_at = ? 
              WHERE id = ?`
	
	_, err := db.db.Exec(ctx, query,
		account.ProviderUsername, account.ProviderEmail, account.ProviderAvatar,
		account.AccessToken, account.RefreshToken, account.TokenExpiresAt,
		account.UpdatedAt, account.ID,
	)
	
	return err
}

// CreateLoginLog 创建登录日志
func (db *DB) CreateLoginLog(ctx context.Context, log *model.LoginLog) error {
	query := `INSERT INTO login_logs (user_id, login_type, ip, user_agent, status, fail_reason, created_at) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.db.Exec(ctx, query,
		log.UserID, log.LoginType, log.IP, log.UserAgent, log.Status, log.FailReason, log.CreatedAt,
	)
	
	return err
}

// GetLoginLogs 获取登录日志（支持分页和筛选）
func (db *DB) GetLoginLogs(ctx context.Context, params map[string]interface{}, page, pageSize int) ([]*model.LoginLog, int, error) {
	// 构建基础查询语句
	baseQuery := `SELECT id, user_id, login_type, ip, user_agent, status, fail_reason, created_at FROM login_logs`
	countQuery := `SELECT COUNT(*) FROM login_logs`
	
	// 构建WHERE条件
	var whereConditions []string
	var args []interface{}
	
	// 根据参数构建筛选条件
	if userID, ok := params["user_id"].(int64); ok && userID > 0 {
		whereConditions = append(whereConditions, "user_id = ?")
		args = append(args, userID)
	}
	
	if ip, ok := params["ip"].(string); ok && ip != "" {
		whereConditions = append(whereConditions, "ip LIKE ?")
		args = append(args, "%"+ip+"%")
	}
	
	if status, ok := params["status"].(int); ok && (status == 0 || status == 1) {
		whereConditions = append(whereConditions, "status = ?")
		args = append(args, status)
	}
	
	if startTime, ok := params["start_time"].(time.Time); ok && !startTime.IsZero() {
		whereConditions = append(whereConditions, "created_at >= ?")
		args = append(args, startTime)
	}
	
	if endTime, ok := params["end_time"].(time.Time); ok && !endTime.IsZero() {
		whereConditions = append(whereConditions, "created_at <= ?")
		args = append(args, endTime)
	}
	
	// 组合WHERE条件
	if len(whereConditions) > 0 {
		whereClause := " WHERE " + strings.Join(whereConditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}
	
	// 添加排序和分页
	baseQuery += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	
	// 添加分页参数
	offset := (page - 1) * pageSize
	args = append(args, pageSize, offset)
	
	// 查询总数
	var total int
	err := db.db.QueryRow(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	
	// 查询数据
	rows, err := db.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	// 解析数据
	var logs []*model.LoginLog
	for rows.Next() {
		log := &model.LoginLog{}
		err := rows.Scan(
			&log.ID, &log.UserID, &log.LoginType, &log.IP,
			&log.UserAgent, &log.Status, &log.FailReason, &log.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}
	
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	
	return logs, total, nil
}

// ListUsers 获取用户列表（支持分页）
func (db *DB) ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int, error) {
	// 计算总数
	var total int
	countQuery := `SELECT COUNT(*) FROM users`
	err := db.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	query := `SELECT id, email, password, username, avatar, role, status, level, email_verified, 
              created_at, updated_at, last_login_at, last_login_ip 
              FROM users ORDER BY id DESC LIMIT ? OFFSET ?`
	
	rows, err := db.db.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.Password, &user.Username, &user.Avatar,
			&user.Role, &user.Status, &user.Level, &user.EmailVerified, &user.CreatedAt,
			&user.UpdatedAt, &user.LastLoginAt, &user.LastLoginIP,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}
	
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	
	return users, total, nil
} 