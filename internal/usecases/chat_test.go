package usecases_test

import (
	"context"
	"testing"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

func TestChat(t *testing.T) {
	testApp := fxtest.New(t,
		coreModules(),
		infrastructureModules(),
		interfaceModules(),
		usecaseModules(),
		fx.Invoke(func(uc *usecases.Chat) {
			got, err := uc.Chat(context.Background(), domain.EventSource{
				Type: linebot.EventSourceTypeUser,
				ID:   "id",
			}, "ドラえもんの誕生日は？")
			if err != nil {
				t.Fatalf("failed to chat: %v", err)
			}
			t.Logf("got: %v", got)
		}),
	)
	testApp.RequireStart().RequireStop()
}
