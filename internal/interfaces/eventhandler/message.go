package eventhandler

import (
	"context"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
	"golang.org/x/text/width"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

type (
	Message struct {
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
) *Message {
	return &Message{
		logger:   logger,
		lineRepo: lineRepo,
		chat:     chat,
		license:  license,
	}
}

func (h *Message) EventType() linebot.EventType {
	return linebot.EventTypeMessage
}

func (h *Message) Handle(ctx context.Context, event *linebot.Event) error {
	if err := h.handle(ctx, event); err != nil {
		h.replyError(event.ReplyToken)
		return err
	}

	return nil
}

func (h *Message) handle(ctx context.Context, event *linebot.Event) error {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		cmd, param, ok := h.extractCmd(message.Text)
		if !ok {
			return nil
		}

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

		switch cmd {
		case commandPrefix:
			err = h.chat.Help(ctx, event.ReplyToken)
		case commandClear:
			err = h.chat.ClearChatHistory(ctx, es, event.ReplyToken)
		case commandApprove:
			err = h.license.Approve(ctx, event.Source.UserID, param)
		case commandReject:
			err = h.license.Reject(ctx, event.Source.UserID, param)
		default:
			err = h.chat.Chat(ctx, es, event.ReplyToken, cmd)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Message) replyError(replyToken string) {
	if err := h.lineRepo.ReplyMessage(context.Background(), replyToken, domain.MessageError); err != nil {
		h.logger.Error("failed to reply message", zap.Error(err))
	}
}

func (h *Message) extractCmd(msg string) (string, string, bool) {
	narrow := strings.Trim(width.Narrow.String(msg), " ")
	if !strings.HasPrefix(narrow, commandPrefix) {
		return "", "", false
	}

	// Help
	if narrow == commandPrefix {
		return commandPrefix, "", true
	}

	narrow = strings.TrimLeft(narrow, commandPrefix)

	// Clear
	if narrow == commandClear {
		return commandClear, "", true
	}

	// Approve or Reject
	split := strings.Split(narrow, " ")
	switch split[0] {
	case commandApprove, commandReject:
		if len(split) == 2 {
			_, err := domain.EventSourceFromUniqueID(split[1])
			if err == nil {
				return split[0], split[1], true
			}

			h.logger.Info("invalid unique key", zap.String("cmd", split[0]), zap.String("unique_key", split[1]))
			return "", "", false
		}
	}

	// AI Chat
	return narrow, "", true
}
