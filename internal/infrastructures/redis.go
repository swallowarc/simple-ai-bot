package infrastructures

import (
	"context"
	"time"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/swallowarc/simple-line-ai-bot/internal/interfaces"
)

const (
	maxRetries = 5
)

type (
	redisClient struct {
		cli redis.Cmdable
	}

	Config struct {
		HostPort string `envconfig:"host_port" default:"localhost:6379"`
		Password string `envconfig:"password" default:""`
		DB       int    `envconfig:"db" default:"0"`
	}
)

func NewRedisClient(config Config) interfaces.MemDBClient {
	return &redisClient{
		cli: redis.NewClient(&redis.Options{
			Addr:       config.HostPort,
			Password:   config.Password, // no password set
			DB:         config.DB,       // use default DB
			MaxRetries: maxRetries,
		}),
	}
}

func (c *redisClient) Ping(ctx context.Context) error {
	if err := c.cli.Ping(ctx).Err(); err != nil {
		return errors.Wrapf(err, "failed to redis Ping")
	}
	return nil
}

func (c *redisClient) Expire(ctx context.Context, key string, duration time.Duration) error {
	if err := c.cli.Expire(ctx, key, duration).Err(); err != nil {
		return errors.Wrap(err, "failed to redis Expire")
	}
	return nil
}

func (c *redisClient) Set(ctx context.Context, key string, value any, duration time.Duration) error {
	if err := c.cli.Set(ctx, key, value, duration).Err(); err != nil {
		return errors.Wrap(err, "failed to redis Set")
	}
	return nil
}

// SetNX sets the value of key to value if key does not exist.
//
//	It returns true if the key was set.
func (c *redisClient) SetNX(ctx context.Context, key string, value any, duration time.Duration) (bool, error) {
	result, err := c.cli.SetNX(ctx, key, value, duration).Result()
	if err != nil {
		return false, errors.Wrap(err, "failed to redis SetNX")
	}
	return result, nil
}

func (c *redisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := c.cli.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", domain.ErrNotFound
	}
	if err != nil {
		return "", errors.Wrap(err, "failed to redis Get")
	}
	return val, nil
}

func (c *redisClient) Del(ctx context.Context, key string) error {
	err := c.cli.Del(ctx, key).Err()
	if err == redis.Nil {
		return domain.ErrNotFound
	}
	if err != nil {
		return errors.Wrap(err, "failed to redis Del")
	}
	return nil
}
