package eventhandler

import (
	"context"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

type (
	Join struct {
		license usecases.License
	}
)

func NewJoinEventHandler(
	license usecases.License,
) *Join {
	return &Join{
		license: license,
	}
}

func (h *Join) EventType() linebot.EventType {
	return linebot.EventTypeJoin
}

func (h *Join) Handle(ctx context.Context, event *linebot.Event) error {
	es, err := convertEventSource(event)
	if err != nil {
		return err
	}

	if _, err := h.license.IssueIfNoLicense(ctx, es, event.ReplyToken); err != nil {
		return err
	}

	return nil
}
