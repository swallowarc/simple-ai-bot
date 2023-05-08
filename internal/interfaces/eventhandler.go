package interfaces

import (
	"context"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
	"golang.org/x/text/width"

	"github.com/swallowarc/lime/lime"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

type (
	messageEventHandler struct {
		logger *zap.Logger
		uc     usecases.Chat
	}
)

func NewMessageEventHandler(logger *zap.Logger, uc usecases.Chat) lime.EventHandler {
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
			h.replyError(cli, event.ReplyToken)
			return err
		}

		cmd, ok := h.extractCmd(message.Text)
		if !ok {
			return nil
		}

		callback := func(ctx context.Context, replyMessage string) error {
			if _, err := cli.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
				h.replyError(cli, event.ReplyToken)
			}
			return nil
		}

		switch cmd {
		case commandPrefix:
			err = h.uc.Help(ctx, callback)
		case commandClear:
			err = h.uc.ClearChatHistory(ctx, es, callback)
		default:
			err = h.uc.Chat(ctx, es, cmd, callback)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *messageEventHandler) replyError(cli lime.LineBotClient, replyToken string) {
	if _, err := cli.ReplyMessage(replyToken, linebot.NewTextMessage(domain.MessageError)).Do(); err != nil {
		h.logger.Error("failed to reply message", zap.Error(err))
	}
}

func (h *messageEventHandler) extractCmd(msg string) (string, bool) {
	narrow := width.Narrow.String(msg)
	nl := len([]rune(narrow))
	pl := len([]rune(commandPrefix))

	if nl < pl || narrow[:pl] != commandPrefix {
		return "", false
	}

	if narrow == commandPrefix {
		return narrow, true
	}

	if nl == pl+1 {
		switch narrow {
		case commandClear:
			return narrow, true
		}
	}

	nr := []rune(narrow)
	cmd := string(nr[pl:])
	return cmd, true
}
