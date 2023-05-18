//go:generate mockgen -source=$GOFILE -destination=../tests/mocks/$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package usecases

import (
	"context"
	"time"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
)

type (
	ChatRepository interface {
		ListCacheMessages(ctx context.Context, es domain.EventSource) (domain.ChatMessages, error)
		UpsertCacheMessages(ctx context.Context, es domain.EventSource, cms domain.ChatMessages) error
		DeleteCacheMessages(ctx context.Context, es domain.EventSource) error

		Chat(ctx context.Context, messages domain.ChatMessages) (domain.ChatMessages, error)
	}

	MessagingRepository interface {
		PushMessages(ctx context.Context, eventSourceID string, messages ...string) error
		ReplyMessages(ctx context.Context, replyToken string, messages ...string) error

		GetGroupName(_ context.Context, groupID string) (string, error)
		ListRoomMemberNames(_ context.Context, roomID string) ([]string, error)
		GetUserName(_ context.Context, userID string) (string, error)
	}

	LicenseRepository interface {
		Get(ctx context.Context, es domain.EventSource) (domain.License, error)
		Upsert(ctx context.Context, lc domain.License, lt time.Duration) error
		Update(ctx context.Context, lc domain.License, lt time.Duration) error
		Delete(ctx context.Context, es domain.EventSource) error
	}
)
