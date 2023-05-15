package usecases_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/fx/fxtest"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
	"github.com/swallowarc/simple-line-ai-bot/internal/interfaces"
	"github.com/swallowarc/simple-line-ai-bot/internal/interfaces/repositories"
	mock_interfaces "github.com/swallowarc/simple-line-ai-bot/internal/tests/mocks/interfaces"
	mock_usecases "github.com/swallowarc/simple-line-ai-bot/internal/tests/mocks/usecases"
	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

func TestChat_Help(t *testing.T) {
	tests := map[string]struct {
		replayToken string
		wantMessage string
	}{
		"help": {
			replayToken: "replyToken",
			wantMessage: domain.MessageHelp(),
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			aiCli := mock_interfaces.NewMockAIClient(ctrl)
			cacheCli := mock_interfaces.NewMockMemDBClient(ctrl)
			msgRepo := mock_usecases.NewMockMessagingRepository(ctrl)

			msgRepo.EXPECT().ReplyMessage(gomock.Any(), tt.replayToken, tt.wantMessage).Return(nil)

			testApp := fxtest.New(t,
				modulesForChat(msgRepo, aiCli, cacheCli),
				fx.Invoke(func(uc usecases.Chat) {
					err := uc.Help(context.Background(), tt.replayToken)
					if err != nil {
						t.Fatalf("failed to help: %v", err)
					}
				}),
			)
			testApp.RequireStart().RequireStop()
		})
	}
}

func TestChat_ClearChatHistory(t *testing.T) {
	tests := map[string]struct {
		es          domain.EventSource
		replayToken string

		cacheMock   func(cache *mock_interfaces.MockMemDBClient)
		aiMock      func(ai *mock_interfaces.MockAIClient)
		msgRepoMock func(repo *mock_usecases.MockMessagingRepository)

		wantErr bool
	}{
		"history not exists": {
			es: domain.EventSource{
				Type: linebot.EventSourceTypeUser,
				ID:   "user-id",
			},
			replayToken: "replyToken",

			cacheMock: func(cache *mock_interfaces.MockMemDBClient) {
				cache.EXPECT().Del(gomock.Any(), "chat_messages:user:user-id").Return(nil)
			},
			msgRepoMock: func(repo *mock_usecases.MockMessagingRepository) {
				repo.EXPECT().ReplyMessage(gomock.Any(), "replyToken", domain.MessageClearHistory).Return(nil)
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			aiCli := mock_interfaces.NewMockAIClient(ctrl)
			if tt.aiMock != nil {
				tt.aiMock(aiCli)
			}
			cacheCli := mock_interfaces.NewMockMemDBClient(ctrl)
			if tt.cacheMock != nil {
				tt.cacheMock(cacheCli)
			}
			msgRepo := mock_usecases.NewMockMessagingRepository(ctrl)
			if tt.msgRepoMock != nil {
				tt.msgRepoMock(msgRepo)
			}

			testApp := fxtest.New(t,
				modulesForChat(msgRepo, aiCli, cacheCli),
				fx.Invoke(func(uc usecases.Chat) {
					err := uc.ClearChatHistory(context.Background(), tt.es, tt.replayToken)
					if err != nil {
						t.Fatalf("failed to help: %v", err)
					}
				}),
			)
			testApp.RequireStart().RequireStop()
		})
	}
}

func modulesForChat(
	msgRepo usecases.MessagingRepository,
	aiCli interfaces.AIClient,
	cacheCli interfaces.MemDBClient,
) fx.Option {
	return fx.Options(
		fx.WithLogger(func() fxevent.Logger {
			return fxevent.NopLogger
		}),
		// interfaces
		fx.Provide(
			repositories.NewChatRepository,
			repositories.NewLicenseRepository,
		),
		// usecases
		fx.Provide(
			usecases.NewChat,
		),
		// mocks
		fx.Provide(
			func() interfaces.AIClient { return aiCli },
			func() interfaces.MemDBClient { return cacheCli },
			func() usecases.MessagingRepository { return msgRepo },
		),
	)
}
