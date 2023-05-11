package eventhandler

import (
	"context"

	"github.com/line/line-bot-sdk-go/v7/linebot"

	"github.com/swallowarc/lime/lime"

	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

type (
	leave struct {
		license usecases.License
	}
)

func NewLeaveHandler(
	license usecases.License,
) lime.EventHandler {
	return &leave{
		license: license,
	}
}

func (h *leave) EventType() linebot.EventType {
	return linebot.EventTypeLeave
}

func (h *leave) Handle(ctx context.Context, event *linebot.Event) error {
	es, err := convertEventSource(event)
	if err != nil {
		return err
	}

	return h.license.Drop(ctx, es)
}
