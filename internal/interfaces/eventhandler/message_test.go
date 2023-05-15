package eventhandler

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
	mock_usecases "github.com/swallowarc/simple-line-ai-bot/internal/tests/mocks/usecases"
)

func TestMessageEventHandler_Handle(t *testing.T) {
	l := zap.NewExample()

	testcases := map[string]struct {
		event      *linebot.Event
		ucModifier func(chat *mock_usecases.MockChat)

		wantErrMsg string
	}{
		"not handled event": {
			event: &linebot.Event{
				Source: &linebot.EventSource{
					Type:   linebot.EventSourceTypeGroup,
					UserID: "test-group-id",
				},
			},
		},
		"success": {
			event: &linebot.Event{
				Source: &linebot.EventSource{
					Type:   linebot.EventSourceTypeUser,
					UserID: "test-user-id",
				},
				Message: &linebot.TextMessage{
					Text: "?test-message",
				},
				ReplyToken: "test-reply-token",
			},
			ucModifier: func(chat *mock_usecases.MockChat) {
				chat.EXPECT().Chat(gomock.Any(), domain.EventSource{
					Type: linebot.EventSourceTypeUser,
					ID:   "test-user-id",
				},
					"test-reply-token",
					"test-message").Return(nil)
			},
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			chat := mock_usecases.NewMockChat(ctrl)
			if tc.ucModifier != nil {
				tc.ucModifier(chat)
			}

			h := &Message{
				logger: l,
				chat:   chat,
			}

			err := h.Handle(nil, tc.event)
			if err != nil {
				if tc.wantErrMsg == "" {
					t.Errorf("unexpected error: %v", err)
				} else if err.Error() != tc.wantErrMsg {
					t.Errorf("error: %v, want error: %v", err, tc.wantErrMsg)
				}
			}
		})
	}
}

func TestMessageEventHandler_extractCmd(t *testing.T) {
	testcases := map[string]struct {
		text string

		wantCmd   string
		wantParam string
		wantOk    bool
	}{
		"not command": {
			text: "test-text",
		},
		"command": {
			text: "?test-text",

			wantCmd:   "test-text",
			wantParam: "",
			wantOk:    true,
		},
		"command help": {
			text: "?",

			wantCmd:   "?",
			wantParam: "",
			wantOk:    true,
		},
		"command help full-width": {
			text: "？",

			wantCmd:   "?",
			wantParam: "",
			wantOk:    true,
		},
		"command clear": {
			text: "?c",

			wantCmd:   "c",
			wantParam: "",
			wantOk:    true,
		},
		"command clear full-width1": {
			text: "？ｃ",

			wantCmd:   "c",
			wantParam: "",
			wantOk:    true,
		},
		"command clear full-width2": {
			text: "？c",

			wantCmd:   "c",
			wantParam: "",
			wantOk:    true,
		},
		"command clear full-width3": {
			text: "?ｃ",

			wantCmd:   "c",
			wantParam: "",
			wantOk:    true,
		},
		"command question": {
			text: "?質問",

			wantCmd:   "質問",
			wantParam: "",
			wantOk:    true,
		},
		"command question full-width1": {
			text: "？質問",

			wantCmd:   "質問",
			wantParam: "",
			wantOk:    true,
		},
		"command question contains command string": {
			text: "?c質問",

			wantCmd:   "c質問",
			wantParam: "",
			wantOk:    true,
		},
		"command approve": {
			text: "?a group:e488b45d-d049-41db-9651-512296cd39db",

			wantCmd:   "a",
			wantParam: "group:e488b45d-d049-41db-9651-512296cd39db",
			wantOk:    true,
		},
		"command  approve(invalid)": {
			text: "?a parameter",

			wantCmd:   "",
			wantParam: "",
			wantOk:    false,
		},
		"command reject": {
			text: "?r group:e488b45d-d049-41db-9651-512296cd39db",

			wantCmd:   "r",
			wantParam: "group:e488b45d-d049-41db-9651-512296cd39db",
			wantOk:    true,
		},
		"command  reject(invalid)": {
			text: "?r parameter",

			wantCmd:   "",
			wantParam: "",
			wantOk:    false,
		},
	}

	h := &Message{
		logger: zap.NewNop(),
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cmd, param, ok := h.extractCmd(tc.text)
			if ok != tc.wantOk {
				t.Errorf("ok: %v, want ok: %v", ok, tc.wantOk)
			}
			if cmd != tc.wantCmd {
				t.Errorf("cmd: %v, want cmd: %v", cmd, tc.wantCmd)
			}
			if param != tc.wantParam {
				t.Errorf("param: %v, want param: %v", param, tc.wantParam)
			}
		})
	}
}
