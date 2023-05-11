package usecases

import (
	"context"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/pkg/errors"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
)

type (
	License interface {
		IssueIfNoLicense(ctx context.Context, es domain.EventSource, replyToken string) (domain.LicenseState, error)
		Approve(ctx context.Context, userID, uniqueKey string) error // Executable only by admin
		Reject(ctx context.Context, userID, uniqueKey string) error  // Executable only by admin
		Drop(ctx context.Context, es domain.EventSource) error
	}

	license struct {
		msgRepo     MessagingRepository
		licenseRepo LicenseRepository

		licenseMode bool
		adminUserID string
	}
)

func NewLicense(
	msgRepo MessagingRepository,
	lineRepo LicenseRepository,
	licenseMode bool,
	adminUserID string,
) License {
	return &license{
		msgRepo:     msgRepo,
		licenseRepo: lineRepo,
		licenseMode: licenseMode,
		adminUserID: adminUserID,
	}
}

func (u *license) IssueIfNoLicense(ctx context.Context, es domain.EventSource, replyToken string) (domain.LicenseState, error) {
	if !u.licenseMode {
		// Always returned as approved if license mode is disabled
		return domain.LicenseStateApproved, nil
	}

	lc, err := u.licenseRepo.Get(ctx, es)
	if errors.Is(err, domain.ErrNotFound) {
		// if license not found, create new license
		lc, err = u.issueLicense(ctx, es)
	}
	if err != nil {
		return domain.LicenseStateNone, err
	}

	if lc.State == domain.LicenseStatePending {
		if err := u.msgRepo.ReplyMessage(ctx, replyToken, domain.MessageLicensePending); err != nil {
			return domain.LicenseStateNone, err
		}
	}

	return lc.State, nil
}

func (u *license) issueLicense(ctx context.Context, es domain.EventSource) (domain.License, error) {
	lc := domain.NewLicense(es)
	err := u.licenseRepo.Upsert(ctx, lc, domain.LicensePendingTTL)
	if err != nil {
		return domain.License{}, err
	}

	var name string
	switch es.Type {
	case linebot.EventSourceTypeGroup:
		name, err = u.msgRepo.GetGroupName(ctx, es.ID)
		if err != nil {
			return domain.License{}, err
		}
	case linebot.EventSourceTypeRoom:
		names, err := u.msgRepo.ListRoomMemberNames(ctx, es.ID)
		if err != nil {
			return domain.License{}, err
		}
		name = strings.Join(names, ",")
	case linebot.EventSourceTypeUser:
		name, err = u.msgRepo.GetUserName(ctx, es.ID)
		if err != nil {
			return domain.License{}, err
		}
	default:
		return domain.License{}, errors.Wrapf(domain.ErrInvalidArgument, "event_source_type: %s", es.Type)
	}

	if err := u.msgRepo.PushMessage(ctx, u.adminUserID, domain.MessageIssueLicense(name, es.ID)); err != nil {
		return domain.License{}, err
	}

	return lc, nil
}

func (u *license) Approve(ctx context.Context, userID, uniqueKey string) error {
	es, err := u.updateLicenseState(ctx, userID, uniqueKey, domain.LicenseStateApproved)
	if err != nil {
		return err
	}

	if err := u.msgRepo.PushMessage(ctx, es.ID, domain.MessageApproved()); err != nil {
		return err
	}

	return nil
}

func (u *license) Reject(ctx context.Context, userID, uniqueKey string) error {
	_, err := u.updateLicenseState(ctx, userID, uniqueKey, domain.LicenseStateRejected)
	if err != nil {
		return err
	}

	// Not send message to user if rejected.
	return nil
}

func (u *license) updateLicenseState(ctx context.Context, userID, uniqueKey string, ls domain.LicenseState) (domain.EventSource, error) {
	if userID != u.adminUserID {
		return domain.EventSource{}, errors.Wrap(domain.ErrPermissionDenied, "failed to update license state")
	}

	es, err := domain.EventSourceFromUniqueID(uniqueKey)
	if err != nil {
		return domain.EventSource{}, err
	}

	lc, err := u.licenseRepo.Get(ctx, es)
	if err != nil {
		return domain.EventSource{}, err
	}

	lc.State = ls
	if err := u.licenseRepo.Update(ctx, lc, 0); err != nil {
		return domain.EventSource{}, err
	}

	return es, nil
}

func (u *license) Drop(ctx context.Context, es domain.EventSource) error {
	if err := u.licenseRepo.Delete(ctx, es); err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			return err
		}
	}

	return nil
}
