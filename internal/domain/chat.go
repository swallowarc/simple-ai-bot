package domain

import (
	"fmt"
	"time"
)

type Role string

const (
	ChatHistoryLimit    = 10
	ChatHistoryLifeHour = 24
	ChatHistoryLife     = 60 * 60 * ChatHistoryLifeHour * time.Second // 1day
)

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

func (r Role) String() string {
	switch r {
	case RoleSystem, RoleUser, RoleAssistant:
		return string(r)
	default:
		return "unknown"
	}
}

type (
	ChatMessage struct {
		Role    Role   `json:"role"`
		Content string `json:"content"`
	}
	ChatMessages []ChatMessage
)

func (cm ChatMessages) LatestMessage() *ChatMessage {
	l := len(cm)
	if l == 0 {
		return nil
	}
	return &cm[l-1]
}

const (
	messageFormatHelp = `【AIの使い方】
  ? : このヘルプを表示
  ?c: AI会話履歴をクリア
  ?(質問文) : AIに質問

過去%d件までの会話履歴を考慮して回答します。
話題を変える場合は会話履歴をクリアしてください。
(最後の会話から%d時間を過ぎると自動クリアされます)`

	MessageClearHistory = "AI会話履歴をクリアしました。"
	MessageError        = "リクエストの処理中にエラーが発生しました。しばらく待ってから再度お試しください。"
)

func MessageHelp() string {
	return fmt.Sprintf(messageFormatHelp, ChatHistoryLimit, ChatHistoryLifeHour)
}
