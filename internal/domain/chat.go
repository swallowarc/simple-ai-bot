package domain

import (
	"time"
)

type Role string

const (
	ChatHistoryLimit   = 10
	ChatHistoryTTLHour = 24
	ChatHistoryTTL     = 60 * 60 * ChatHistoryTTLHour * time.Second // 1day
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
