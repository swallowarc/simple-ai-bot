package eventhandler

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
)

const (
	commandPrefix = "?"
	commandClear  = commandPrefix + "c"
)

func convertEventSource(event *linebot.Event) (domain.EventSource, error) {
	var id string
	switch event.Source.Type {
	case linebot.EventSourceTypeUser:
		id = event.Source.UserID
	case linebot.EventSourceTypeGroup:
		id = event.Source.GroupID
	case linebot.EventSourceTypeRoom:
		id = event.Source.RoomID
	default:
		return domain.EventSource{}, domain.ErrUnknownEventSource
	}

	return domain.EventSource{
		Type: event.Source.Type,
		ID:   id,
	}, nil
}
