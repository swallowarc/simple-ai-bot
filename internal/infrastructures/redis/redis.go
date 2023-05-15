package redis

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/go-redis/redis/v8"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
	"github.com/swallowarc/simple-line-ai-bot/internal/interfaces"
)

const (
	maxRetries = 5
)

type (
	client struct {
		cli redis.Cmdable
	}

	Env struct {
		HostPort string `envconfig:"host_port" default:"localhost:6379"`
		Password string `envconfig:"password" default:""`
		DB       int    `envconfig:"db" default:"0"`
	}
)

func NewClient(env Env) interfaces.MemDBClient {
	return &client{
		cli: redis.NewClient(&redis.Options{
			Addr:       env.HostPort,
			Password:   env.Password,
			DB:         env.DB,
			MaxRetries: maxRetries,
		}),
	}
}

func (c *client) Ping(ctx context.Context) error {
	if err := c.cli.Ping(ctx).Err(); err != nil {
		return errors.Wrapf(err, "failed to redis Ping")
	}
	return nil
}

func (c *client) Expire(ctx context.Context, key string, duration time.Duration) error {
	if err := c.cli.Expire(ctx, key, duration).Err(); err != nil {
		return errors.Wrap(err, "failed to redis Expire")
	}
	return nil
}

func (c *client) Set(ctx context.Context, key string, value any, duration time.Duration) error {
	if err := c.cli.Set(ctx, key, value, duration).Err(); err != nil {
		return errors.Wrap(err, "failed to redis Set")
	}
	return nil
}

// SetNX sets the value of key to value if key does not exist.
//
//	It returns true if the key was set.
func (c *client) SetNX(ctx context.Context, key string, value any, duration time.Duration) (bool, error) {
	result, err := c.cli.SetNX(ctx, key, value, duration).Result()
	if err != nil {
		return false, errors.Wrap(err, "failed to redis SetNX")
	}
	return result, nil
}

// SetXX sets the value of key to value if key already exists.
//
//	It returns true if the key was set.
func (c *client) SetXX(ctx context.Context, key string, value any, duration time.Duration) (bool, error) {
	result, err := c.cli.SetXX(ctx, key, value, duration).Result()
	if err != nil {
		return false, errors.Wrap(err, "failed to redis SetXX")
	}
	return result, nil
}

func (c *client) Get(ctx context.Context, key string) (string, error) {
	val, err := c.cli.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", domain.ErrNotFound
	}
	if err != nil {
		return "", errors.Wrap(err, "failed to redis Get")
	}
	return val, nil
}

func (c *client) Del(ctx context.Context, key string) error {
	result, err := c.cli.Del(ctx, key).Result()
	if err != nil {
		return errors.Wrap(err, "failed to redis Del")
	}
	if result == 0 {
		// noop: not error if key does not exist
		// return domain.ErrNotFound
	}
	return nil
}
