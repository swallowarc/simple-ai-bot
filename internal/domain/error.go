package domain

import "github.com/pkg/errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrUnknownEventSource = errors.New("unknown event source")
)
