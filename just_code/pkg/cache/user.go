package cache

import (
	"context"
	"encoding/json"
	"golang-jwttoken/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
)

const UserCacheTimeout = 10 * time.Minute

type User interface {
	Get(ctx context.Context, key string) (*model.Users, error)
	Set(ctx context.Context, key string, value *model.Users) error
}

type UserCache struct {
	Expiration time.Duration
	redisCli   *redis.Client
}

func NewUserCache(redisCli *redis.Client, expiration time.Duration) User {
	return &UserCache{
		redisCli:   redisCli,
		Expiration: expiration,
	}
}

func (c *UserCache) Get(ctx context.Context, key string) (*model.Users, error) {
	value := c.redisCli.Get(ctx, key).Val()

	if value == "" {
		return nil, nil
	}

	var user *model.Users
	err := json.Unmarshal([]byte(value), &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c *UserCache) Set(ctx context.Context, key string, value *model.Users) error {
	userJson, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.redisCli.Set(ctx, key, string(userJson), c.Expiration).Err()
}
