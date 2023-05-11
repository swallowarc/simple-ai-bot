package eventhandler

import (
	"context"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/swallowarc/lime/lime"

	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

type (
	join struct {
		license usecases.License
	}
)

func NewJoinEventHandler(
	license usecases.License,
) lime.EventHandler {
	return &join{
		license: license,
	}
}

func (h *join) EventType() linebot.EventType {
	return linebot.EventTypeJoin
}

func (h *join) Handle(ctx context.Context, event *linebot.Event) error {
	es, err := convertEventSource(event)
	if err != nil {
		return err
	}

	if _, err := h.license.IssueIfNoLicense(ctx, es, event.ReplyToken); err != nil {
		return err
	}

	return nil
}
