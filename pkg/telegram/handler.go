package telegram

import (
	"context"
	"encoding/base64"
	"github.com/Zkeai/DDPay/internal/model"
	"github.com/Zkeai/DDPay/pkg/redis"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// decodeBase64URL 将 URL 编码的 Base64 字符串转换为标准 Base64，并进行解码
func decodeBase64URL(base64url string) ([]byte, error) {
	// 替换 URL 编码的 Base64 字符为标准 Base64 字符
	base64str := strings.ReplaceAll(base64url, "-", "+")
	base64str = strings.ReplaceAll(base64str, "_", "/")

	// 添加填充字符 "="，确保长度是4的倍数
	padding := len(base64str) % 4
	if padding > 0 {
		base64str += strings.Repeat("=", 4-padding)
	}

	// 解码
	return base64.StdEncoding.DecodeString(base64str)
}

func (s *TelegramService) StartHandler() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := s.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		// 机器人被拉进群组
		if update.Message.NewChatMembers != nil {
			for _, member := range update.Message.NewChatMembers {
				if member.UserName == s.Bot.Self.UserName {
					chatID := update.Message.Chat.ID
					chatTitle := update.Message.Chat.Title
					msg := tgbotapi.NewMessage(chatID, "Hi there,\n\nYour group name is: "+chatTitle+"\n\nYour chat id is: "+strconv.FormatInt(chatID, 10))
					_, _ = s.Bot.Send(msg)
				}
			}
		}

		// 监听 /start 带参数
		if update.Message.IsCommand() && update.Message.Command() == "start" {
			params := update.Message.CommandArguments() // 获取 start 后的参数
			s.handleStartCommand(update.Message.Chat.ID, params)
		}
	}
}

func (s *TelegramService) handleStartCommand(chatID int64, args string) {

	// 使用自定义的 URL Base64 解码函数
	data, err := decodeBase64URL(args)
	if err != nil {
		err := redis.Set(args, "failed", 6*time.Minute)
		if err != nil {
			return
		}
		return
	}

	// 将解码后的字节数据转换为字符串
	decodedString := string(data)

	// 解析解码后的数据
	parts := strings.Split(decodedString, "|")
	if len(parts) != 4 {
		err := redis.Set(args, "failed", 6*time.Minute)
		if err != nil {
			return
		}

		return
	}

	platform, userID, _, suffix := parts[0], parts[1], parts[2], parts[3]

	// 校验平台和后缀
	if platform != "T" || suffix != "U" {
		err := redis.Set(args, "failed", 6*time.Minute)
		if err != nil {
			return
		}
		return
	}
	tgName, err := redis.Get("tgName")
	if err != nil {
		return
	}
	// 返回成功信息

	userIDInt, _ := strconv.Atoi(userID)
	var req *model.ChannelDTO
	req = &model.ChannelDTO{
		UserID:    userIDInt,
		ChannelID: "telegram",
		Name:      tgName,
		ChatID:    strconv.FormatInt(chatID, 10),
	}
	err = CallUpsertChannel(context.Background(), req)
	if err != nil {
		return
	}

	//发送消息到bot
	msg := tgbotapi.NewMessage(chatID, "✅ AIDOG Telegram通知渠道创建")
	_, _ = s.Bot.Send(msg)
	//redis 置成功状态
	_ = redis.Set(args, "success", 10*time.Minute)

}
