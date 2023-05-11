package domain

import (
	"fmt"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/pkg/errors"
)

type EventSource struct {
	Type linebot.EventSourceType
	ID   string
}

func (es EventSource) UniqueID() string {
	return fmt.Sprintf("%s:%s", es.Type, es.ID)
}

func EventSourceFromUniqueID(uniqueID string) (EventSource, error) {
	sep := strings.Split(uniqueID, ":")
	if len(sep) != 2 {
		return EventSource{}, errors.Errorf("invalid unique id: %s", uniqueID)
	}

	est := linebot.EventSourceType(sep[0])
	switch est {
	case linebot.EventSourceTypeGroup, linebot.EventSourceTypeRoom, linebot.EventSourceTypeUser:
		// noop
	default:
		return EventSource{}, errors.Errorf("invalid unique id: %s", uniqueID)
	}

	return EventSource{
		Type: est,
		ID:   sep[1],
	}, nil
}
