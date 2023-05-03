package domain

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type (
	EventSource struct {
		Type linebot.EventSourceType
		ID   string
	}
)

func (e EventSource) Key() string {
	return fmt.Sprintf("%s:%s", e.Type, e.ID)
}
