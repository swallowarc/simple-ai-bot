package interfaces

import (
	"context"

	"github.com/line/line-bot-sdk-go/v7/linebot"

	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

type (
	messageEventHandler struct {
		uc *usecases.Chat
	}
)

func (m messageEventHandler) Handle(ctx context.Context, event *linebot.Event, cli *linebot.Client) error {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		if _, err := cli.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
			event.
		}
	}

	return nil
}
