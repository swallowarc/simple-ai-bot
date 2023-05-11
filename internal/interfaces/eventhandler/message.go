package eventhandler

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

		lineRepo usecases.MessagingRepository
		chat     usecases.Chat
		license  usecases.License
	}
)

func NewMessageEventHandler(
	logger *zap.Logger,
	lineRepo usecases.MessagingRepository,
	chat usecases.Chat,
	license usecases.License,
) lime.EventHandler {
	return &messageEventHandler{
		logger:   logger,
		lineRepo: lineRepo,
		chat:     chat,
		license:  license,
	}
}

func (h *messageEventHandler) EventType() linebot.EventType {
	return linebot.EventTypeMessage
}

func (h *messageEventHandler) Handle(ctx context.Context, event *linebot.Event) error {
	if err := h.handle(ctx, event); err != nil {
		h.replyError(event.ReplyToken)
		return err
	}

	return nil
}

func (h *messageEventHandler) handle(ctx context.Context, event *linebot.Event) error {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		es, err := convertEventSource(event)
		if err != nil {
			return err
		}

		// Check license
		ls, err := h.license.IssueIfNoLicense(ctx, es, event.ReplyToken)
		if err != nil {
			return err
		}
		if !ls.IsApproved() {
			h.logger.Info("license is not approved",
				zap.String("unique_key", es.UniqueID()), zap.String("state", ls.String()))
			return nil
		}

		cmd, ok := h.extractCmd(message.Text)
		if !ok {
			return nil
		}

		switch cmd {
		case commandPrefix:
			err = h.chat.Help(ctx, event.ReplyToken)
		case commandClear:
			err = h.chat.ClearChatHistory(ctx, es, event.ReplyToken)
		// TODO: approveとrejectの処理を実装
		default:
			err = h.chat.Chat(ctx, es, event.ReplyToken, cmd)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *messageEventHandler) replyError(replyToken string) {
	if err := h.lineRepo.ReplyMessage(context.Background(), replyToken, domain.MessageError); err != nil {
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
