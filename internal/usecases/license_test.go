package usecases_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/fx/fxtest"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
	"github.com/swallowarc/simple-line-ai-bot/internal/infrastructures/env"
	"github.com/swallowarc/simple-line-ai-bot/internal/infrastructures/redis"
	"github.com/swallowarc/simple-line-ai-bot/internal/interfaces"
	"github.com/swallowarc/simple-line-ai-bot/internal/interfaces/repositories"
	mock_usecases "github.com/swallowarc/simple-line-ai-bot/internal/tests/mocks/usecases"
	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

func TestLicense_IssueIfNoLicense(t *testing.T) {
	adminUserID := uuid.NewString()

	tests := map[string]struct {
		// env
		licenseMode bool

		// input
		es          domain.EventSource
		replayToken string

		// mocks
		msgRepoMock func(repo *mock_usecases.MockMessagingRepository, userID string)
		init        func(ctx context.Context, memDBCli interfaces.MemDBClient) error

		// want
		wantState domain.LicenseState
		wantErr   bool
	}{
		"license mode off": {
			licenseMode: false,

			es: domain.EventSource{
				Type: linebot.EventSourceTypeUser,
				ID:   uuid.NewString(),
			},
			replayToken: "replyToken",

			wantState: domain.LicenseStateApproved,
		},
		"license none": {
			licenseMode: true,

			es: domain.EventSource{
				Type: linebot.EventSourceTypeUser,
				ID:   uuid.NewString(),
			},
			replayToken: "replyToken",

			msgRepoMock: func(repo *mock_usecases.MockMessagingRepository, id string) {
				repo.EXPECT().ReplyMessage(gomock.Any(), "replyToken", domain.MessageLicensePending).Return(nil)
				repo.EXPECT().GetUserName(gomock.Any(), id).Return("user-name", nil)
				repo.EXPECT().PushMessage(gomock.Any(), adminUserID, domain.MessageIssueLicense("user", "user-name", id)).Return(nil)
			},

			wantState: domain.LicenseStatePending,
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			msgRepo := mock_usecases.NewMockMessagingRepository(ctrl)
			if tt.msgRepoMock != nil {
				tt.msgRepoMock(msgRepo, tt.es.ID)
			}

			testApp := fxtest.New(t,
				// usecases
				fx.Provide(
					func(msgRepo usecases.MessagingRepository, licenseRepo usecases.LicenseRepository) usecases.License {
						return usecases.NewLicense(msgRepo, licenseRepo, tt.licenseMode, adminUserID)
					},
				),
				modulesForLicense(msgRepo),
				fx.Invoke(
					func(lc fx.Lifecycle, memDBCli interfaces.MemDBClient) {
						if tt.init != nil {
							lc.Append(fx.Hook{
								OnStart: func(context.Context) error {
									return tt.init(context.Background(), memDBCli)
								},
							})
						}
					},
					func(uc usecases.License) {
						state, err := uc.IssueIfNoLicense(context.Background(), tt.es, tt.replayToken)
						if (err != nil) != tt.wantErr {
							t.Fatalf("want err: %v, but got: %v", tt.wantErr, err)
						}
						if state != tt.wantState {
							t.Errorf("want state: %v, but got: %v", tt.wantState, state)
						}
					},
				),
			)
			testApp.RequireStart().RequireStop()
		})
	}
}

func modulesForLicense(
	msgRepo usecases.MessagingRepository,
) fx.Option {
	return fx.Options(
		fx.WithLogger(func() fxevent.Logger {
			return fxevent.NopLogger
		}),
		// infrastructures
		fx.Provide(
			env.GetRedisEnv,
			redis.NewClient,
		),
		// interfaces
		fx.Provide(
			repositories.NewLicenseRepository,
		),
		// mocks
		fx.Provide(
			func() usecases.MessagingRepository { return msgRepo },
		),
	)
}
