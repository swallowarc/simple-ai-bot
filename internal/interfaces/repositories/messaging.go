package repositories

import (
	"context"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/pkg/errors"

	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

type (
	messagingRepository struct {
		cli *linebot.Client
	}
)

func NewMessagingRepository(cli *linebot.Client) usecases.MessagingRepository {
	return &messagingRepository{
		cli: cli,
	}
}

func (r *messagingRepository) PushMessage(_ context.Context, eventSourceID string, message string) error {
	if _, err := r.cli.PushMessage(eventSourceID, linebot.NewTextMessage(message)).Do(); err != nil {
		return errors.Wrap(err, "failed to push message")
	}

	return nil
}

func (r *messagingRepository) ReplyMessage(_ context.Context, replyToken string, message string) error {
	if _, err := r.cli.ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do(); err != nil {
		return errors.Wrap(err, "failed to reply message")
	}

	return nil
}

func (r *messagingRepository) GetGroupName(_ context.Context, groupID string) (string, error) {
	s, err := r.cli.GetGroupSummary(groupID).Do()
	if err != nil {
		return "", errors.Wrap(err, "failed to get group summary")
	}

	return s.GroupName, nil
}

func (r *messagingRepository) ListRoomMemberNames(_ context.Context, roomID string) ([]string, error) {
	names := make([]string, 0)
	ct := ""
	for {
		members, err := r.cli.GetRoomMemberIDs(roomID, ct).Do()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get room members")
		}

		for _, m := range members.MemberIDs {
			p, err := r.cli.GetProfile(m).Do()
			if err != nil {
				return nil, errors.Wrap(err, "failed to get profile")
			}
			names = append(names, p.DisplayName)
		}

		if members.Next != "" {
			ct = members.Next
		} else {
			break
		}
	}

	return names, nil
}

func (r *messagingRepository) GetUserName(_ context.Context, userID string) (string, error) {
	p, err := r.cli.GetProfile(userID).Do()
	if err != nil {
		return "", errors.Wrap(err, "failed to get profile")
	}

	return p.DisplayName, nil
}
