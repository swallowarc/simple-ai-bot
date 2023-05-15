package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
	"github.com/swallowarc/simple-line-ai-bot/internal/interfaces"
	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

type (
	licenseRepository struct {
		memDBCli interfaces.MemDBClient
	}
)

const (
	keyLicense = "license:%s"
)

func (r *licenseRepository) licenseKey(es domain.EventSource) string {
	return fmt.Sprintf(keyLicense, es.UniqueID())
}

func (r *licenseRepository) Get(ctx context.Context, es domain.EventSource) (domain.License, error) {
	lic, err := getFromMemDB[domain.License](ctx, r.memDBCli, r.licenseKey(es))
	if err != nil {
		return domain.License{}, errors.Wrapf(err, "failed to get license")
	}
	return lic, nil
}

func (r *licenseRepository) Upsert(ctx context.Context, lc domain.License, lt time.Duration) error {
	return setToMemDB[domain.License](ctx, r.memDBCli, r.licenseKey(lc.EventSource), lc, lt)
}

func (r *licenseRepository) Update(ctx context.Context, lc domain.License, lt time.Duration) error {
	updated, err := setXXToMemDB(ctx, r.memDBCli, r.licenseKey(lc.EventSource), lc, lt)
	if err != nil {
		return err
	}
	if !updated {
		return errors.Wrap(domain.ErrNotFound, "failed to update license")
	}
	return nil
}

func (r *licenseRepository) Delete(ctx context.Context, es domain.EventSource) error {
	if err := r.memDBCli.Del(ctx, r.licenseKey(es)); err != nil {
		return errors.Wrapf(err, "failed to delete")
	}
	return nil
}

func NewLicenseRepository(memDBCli interfaces.MemDBClient) usecases.LicenseRepository {
	return &licenseRepository{
		memDBCli: memDBCli,
	}
}
