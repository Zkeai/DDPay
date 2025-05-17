package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type UserChannel struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	ChannelID      string    `json:"channel_id"`
	WebhookURL     string    `json:"webhook_url"`
	Channel        string    `json:"channel"`
	Username       string    `json:"username"`
	Name           string    `json:"name"`
	ChatID         string    `json:"chat_id"`
	EmailAddress   string    `json:"email_address"`
	RegID          string    `json:"reg_id"`
	LastSend       int64     `json:"last_send"`
	FailedStart    int       `json:"failed_start"`
	FailedCount    int       `json:"failed_count"`
	FailedStage    int       `json:"failed_stage"`
	FailedAvail    int       `json:"failed_avail"`
	Disabled       int       `json:"disabled"`
	DisabledReason string    `json:"disabled_reason"`
	IsDefault      int       `json:"is_default"`
	TelBotTokens   string    `json:"tel_bot_tokens"`
	IsDelete       int       `json:"is_delete"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Exists 判断是否已存在该用户和通道
func (db *DB) Exists(ctx context.Context, userID int, channelID string) (bool, error) {
	var count int
	err := db.db.QueryRowContext(ctx, `
		SELECT COUNT(1) FROM yu_channel 
		WHERE user_id = ? AND channel_id = ? AND is_delete = 0
	`, userID, channelID).Scan(&count)
	return count > 0, err
}

// Create 创建新的通道
func (db *DB) Create(ctx context.Context, c *UserChannel) error {
	query := `
		INSERT INTO yu_channel (
			user_id, channel_id, webhook_url, channel, username, name,
			chat_id, email_address, reg_id, tel_bot_tokens, is_default,
			disabled, is_delete, last_send, failed_start, failed_count, 
			failed_stage, failed_avail, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, 0, 0, 0, 0, 0, 0, NOW(), NOW())
	`

	_, err := db.db.ExecContext(ctx, query,
		c.UserID, c.ChannelID, c.WebhookURL, c.Channel, c.Username, c.Name,
		c.ChatID, c.EmailAddress, c.RegID, c.TelBotTokens, c.IsDefault,
	)
	return err
}

// Update 更新已有通道
func (db *DB) Update(ctx context.Context, userID int, channelID string, c *UserChannel) error {
	var fields []string
	var args []interface{}

	if c.WebhookURL != "" {
		fields = append(fields, "webhook_url = ?")
		args = append(args, c.WebhookURL)
	}
	if c.Channel != "" {
		fields = append(fields, "channel = ?")
		args = append(args, c.Channel)
	}
	if c.Username != "" {
		fields = append(fields, "username = ?")
		args = append(args, c.Username)
	}
	if c.Name != "" {
		fields = append(fields, "name = ?")
		args = append(args, c.Name)
	}
	if c.ChatID != "" {
		fields = append(fields, "chat_id = ?")
		args = append(args, c.ChatID)
	}
	if c.EmailAddress != "" {
		fields = append(fields, "email_address = ?")
		args = append(args, c.EmailAddress)
	}
	if c.RegID != "" {
		fields = append(fields, "reg_id = ?")
		args = append(args, c.RegID)
	}
	if c.TelBotTokens != "" {
		fields = append(fields, "tel_bot_tokens = ?")
		args = append(args, c.TelBotTokens)
	}
	fields = append(fields, "is_default = ?")
	args = append(args, c.IsDefault)

	fields = append(fields, "updated_at = ?")
	args = append(args, time.Now())

	args = append(args, userID, channelID)

	query := fmt.Sprintf(`
		UPDATE yu_channel SET %s 
		WHERE user_id = ? AND channel_id = ? AND is_delete = 0
	`, strings.Join(fields, ", "))

	_, err := db.db.ExecContext(ctx, query, args...)
	return err
}

// GetByUserID 查询用户所有通道（未删除）
func (db *DB) GetByUserID(ctx context.Context, userID int) ([]*UserChannel, error) {
	query := `
		SELECT id, user_id, channel_id, webhook_url, channel, username, name,
		       chat_id, email_address, reg_id, last_send, failed_start, failed_count,
		       failed_stage, failed_avail, disabled, disabled_reason, is_default,
		       tel_bot_tokens, is_delete, created_at, updated_at
		FROM yu_channel 
		WHERE user_id = ? AND is_delete = 0
	`

	rows, err := db.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var channels []*UserChannel
	for rows.Next() {
		var c UserChannel
		err := rows.Scan(
			&c.ID, &c.UserID, &c.ChannelID, &c.WebhookURL, &c.Channel, &c.Username, &c.Name,
			&c.ChatID, &c.EmailAddress, &c.RegID, &c.LastSend, &c.FailedStart, &c.FailedCount,
			&c.FailedStage, &c.FailedAvail, &c.Disabled, &c.DisabledReason, &c.IsDefault,
			&c.TelBotTokens, &c.IsDelete, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		channels = append(channels, &c)
	}
	return channels, nil
}

// SetDisabled 启用 / 禁用通道
func (db *DB) SetDisabled(ctx context.Context, userID int, channelID string, disabled int, reson string) error {
	query := `
		UPDATE yu_channel SET disabled = ?,disabled_reason = ?, updated_at = NOW()
		WHERE user_id = ? AND channel_id = ? AND is_delete = 0
	`
	_, err := db.db.ExecContext(ctx, query, disabled, reson, userID, channelID)
	return err
}

// SoftDelete 软删除通道
func (db *DB) SoftDelete(ctx context.Context, userID int, channelID string) error {
	query := `
		UPDATE yu_channel SET is_delete = 1, updated_at = NOW()
		WHERE user_id = ? AND channel_id = ? AND is_delete = 0
	`
	_, err := db.db.ExecContext(ctx, query, userID, channelID)
	return err
}
