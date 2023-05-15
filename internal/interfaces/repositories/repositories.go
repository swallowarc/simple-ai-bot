package repositories

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/swallowarc/simple-line-ai-bot/internal/interfaces"
)

func toJson[T any](t T) (string, error) {
	j, err := json.Marshal(t)
	if err != nil {
		return "", errors.New("failed to marshal json")
	}
	return string(j), nil
}

func setToMemDB[T any](ctx context.Context, memDB interfaces.MemDBClient, key string, t T, d time.Duration) error {
	j, err := toJson[T](t)
	if err != nil {
		return err
	}
	return memDB.Set(ctx, key, j, d)
}

func setXXToMemDB[T any](ctx context.Context, memDB interfaces.MemDBClient, key string, t T, d time.Duration) (bool, error) {
	j, err := toJson[T](t)
	if err != nil {
		return false, err
	}
	return memDB.SetXX(ctx, key, j, d)
}

func fromJson[T any](j string) (T, error) {
	var t T
	if err := json.Unmarshal([]byte(j), &t); err != nil {
		return t, errors.New("failed to unmarshal json")
	}
	return t, nil
}

func getFromMemDB[T any](ctx context.Context, memDB interfaces.MemDBClient, key string) (T, error) {
	jsn, err := memDB.Get(ctx, key)
	if err != nil {
		var t T
		return t, err
	}
	return fromJson[T](jsn)
}
