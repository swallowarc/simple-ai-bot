package eventhandler

import (
	"context"

	"github.com/line/line-bot-sdk-go/v7/linebot"

	"github.com/swallowarc/lime/lime"

	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

type (
	unfollow struct {
		license usecases.License
	}
)

func NewUnfollowHandler(
	license usecases.License,
) lime.EventHandler {
	return &unfollow{
		license: license,
	}
}

func (h *unfollow) EventType() linebot.EventType {
	return linebot.EventTypeUnfollow
}

func (h *unfollow) Handle(ctx context.Context, event *linebot.Event) error {
	es, err := convertEventSource(event)
	if err != nil {
		return err
	}

	return h.license.Drop(ctx, es)
}
