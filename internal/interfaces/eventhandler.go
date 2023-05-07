package interfaces

import (
	"context"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/swallowarc/lime/lime"
	"go.uber.org/zap"
	"golang.org/x/text/width"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

type (
	messageEventHandler struct {
		logger *zap.Logger
		uc     *usecases.Chat
	}
)

const (
	commandPrefix = "?"
	commandClear  = commandPrefix + "c"

	helpText = `AIの使い方
* ?: このヘルプを表示
* ?c: AI会話履歴をクリア
* ?(質問文): AIに質問
`
)

func NewMessageEventHandler(logger *zap.Logger, uc *usecases.Chat) lime.EventHandler {
	return &messageEventHandler{
		logger: logger,
		uc:     uc,
	}
}

func convEventMessage(event *linebot.Event) (domain.EventSource, error) {
	var id string
	switch event.Source.Type {
	case linebot.EventSourceTypeUser:
		id = event.Source.UserID
	case linebot.EventSourceTypeGroup:
		id = event.Source.GroupID
	case linebot.EventSourceTypeRoom:
		id = event.Source.RoomID
	default:
		return domain.EventSource{}, domain.ErrUnknownEventSource
	}

	return domain.EventSource{
		Type: event.Source.Type,
		ID:   id,
	}, nil
}

func (h *messageEventHandler) Handle(ctx context.Context, event *linebot.Event, cli lime.LineBotClient) error {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		es, err := convEventMessage(event)
		if err != nil {
			return err
		}

		prefix := width.Narrow.String(message.Text[:len(commandPrefix)])
		if prefix != commandPrefix {
			return nil
		}

		switch {
		case message.Text == commandPrefix:
			if _, err := cli.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(helpText)).Do(); err != nil {
				return err
			}

		case message.Text == commandClear:
			if err := h.uc.ClearChatHistory(ctx, es); err != nil {
				return err
			}

			if _, err := cli.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("AI会話履歴をクリアしました")).Do(); err != nil {
				return err
			}

		default:
			aiResponse, err := h.uc.Chat(ctx, es, message.Text[len(commandPrefix):])
			if err != nil {
				return err
			}

			if _, err := cli.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(aiResponse)).Do(); err != nil {
				return err
			}
		}
	}

	return nil
}
