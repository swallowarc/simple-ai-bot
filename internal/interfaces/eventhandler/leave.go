package eventhandler

import (
	"context"

	"github.com/line/line-bot-sdk-go/v7/linebot"

	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

type (
	Leave struct {
		license usecases.License
	}
)

func NewLeaveHandler(
	license usecases.License,
) *Leave {
	return &Leave{
		license: license,
	}
}

func (h *Leave) EventType() linebot.EventType {
	return linebot.EventTypeLeave
}

func (h *Leave) Handle(ctx context.Context, event *linebot.Event) error {
	es, err := convertEventSource(event)
	if err != nil {
		return err
	}

	return h.license.Drop(ctx, es)
}
