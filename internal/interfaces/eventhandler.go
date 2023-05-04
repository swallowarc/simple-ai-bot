package interfaces

import (
	"context"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/swallowarc/lime/lime"
	"go.uber.org/zap"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

type (
	messageEventHandler struct {
		logger *zap.Logger
		uc     *usecases.Chat
	}
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

func (h *messageEventHandler) Handle(ctx context.Context, event *linebot.Event, cli *linebot.Client) error {
	if err := h.handleLogic(ctx, event, cli); err != nil {
		h.logger.Error("failed to handle logic", zap.Error(err))
		return err
	}

	return nil
}

func (h *messageEventHandler) handleLogic(ctx context.Context, event *linebot.Event, cli *linebot.Client) error {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		if _, err := cli.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
			es, err := convEventMessage(event)
			if err != nil {
				return err
			}

			aiResponse, err := h.uc.Chat(ctx, es, message.Text)
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
