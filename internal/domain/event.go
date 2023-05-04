package domain

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type EventSource struct {
	Type linebot.EventSourceType
	ID   string
}
