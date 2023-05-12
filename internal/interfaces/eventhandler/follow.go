package eventhandler

import (
	"context"

	"github.com/line/line-bot-sdk-go/v7/linebot"

	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

type (
	Follow struct {
		license usecases.License
	}
)

func NewFollowHandler(
	license usecases.License,
) *Follow {
	return &Follow{
		license: license,
	}
}

func (h *Follow) EventType() linebot.EventType {
	return linebot.EventTypeUnfollow
}

func (h *Follow) Handle(ctx context.Context, event *linebot.Event) error {
	es, err := convertEventSource(event)
	if err != nil {
		return err
	}

	if _, err := h.license.IssueIfNoLicense(ctx, es, event.ReplyToken); err != nil {
		return err
	}

	return nil
}
