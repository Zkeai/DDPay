package dto

type ChannelDTO struct {
	UserID         int    `json:"user_id" binding:"required"`
	ChannelID      string `json:"channel_id" binding:"required"`
	WebhookURL     string `json:"webhook_url"`
	Channel        string `json:"channel"`
	Username       string `json:"username"`
	Name           string `json:"name"`
	ChatID         string `json:"chat_id"`
	EmailAddress   string `json:"email_address"`
	RegID          string `json:"reg_id"`
	LastSend       int64  `json:"last_send"`
	FailedStart    int    `json:"failed_start"`
	FailedCount    int    `json:"failed_count"`
	FailedStage    int    `json:"failed_stage"`
	FailedAvail    int    `json:"failed_avail"`
	Disabled       int    `json:"disabled"`
	DisabledReason string `json:"disabled_reason"`
	IsDefault      int    `json:"is_default"`
	TelBotTokens   string `json:"tel_bot_tokens"`
}

type ChannelStatusUpdateDTO struct {
	UserID    int    `json:"user_id" binding:"required"`
	ChannelID string `json:"channel_id" binding:"required"`
	Disabled  int    `json:"disabled"` // 0: 启用, 1: 禁用
}

type ChannelDeleteDTO struct {
	UserID    int    `json:"user_id" binding:"required"`
	ChannelID string `json:"channel_id" binding:"required"`
}
