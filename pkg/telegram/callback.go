// common/telegram/callback.go

package telegram

import (
	"context"

	"github.com/Zkeai/DDPay/internal/model"
)

var upsertChannelFunc func(ctx context.Context, dto *model.ChannelDTO) error

// RegisterUpsertChannelHandler 注入上层处理逻辑（例如 service.UpsertChannel）
func RegisterUpsertChannelHandler(handler func(ctx context.Context, dto *model.ChannelDTO) error) {
	upsertChannelFunc = handler
}

func CallUpsertChannel(ctx context.Context, dto *model.ChannelDTO) error {
	if upsertChannelFunc == nil {
		return nil
	}
	return upsertChannelFunc(ctx, dto)
}
