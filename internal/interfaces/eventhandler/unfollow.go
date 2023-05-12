package eventhandler

import (
	"context"

	"github.com/line/line-bot-sdk-go/v7/linebot"

	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

type (
	Unfollow struct {
		license usecases.License
	}
)

func NewUnfollowHandler(
	license usecases.License,
) *Unfollow {
	return &Unfollow{
		license: license,
	}
}

func (h *Unfollow) EventType() linebot.EventType {
	return linebot.EventTypeUnfollow
}

func (h *Unfollow) Handle(ctx context.Context, event *linebot.Event) error {
	es, err := convertEventSource(event)
	if err != nil {
		return err
	}

	return h.license.Drop(ctx, es)
}
