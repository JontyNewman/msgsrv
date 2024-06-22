package redis

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisMessageRepository struct {
	ctx    context.Context
	client *redis.Client
}

func NewRedisMessageRepository(ctx context.Context, client *redis.Client) RedisMessageRepository {
	return RedisMessageRepository{ctx: ctx, client: client}
}

func (repo *RedisMessageRepository) Add(message string) (uint, error) {
	client := repo.client
	ctx := repo.ctx

	incr := client.Incr(ctx, "next_message_id")

	if err := incr.Err(); err != nil {
		return 0, err
	}

	id := incr.Val()

	set := client.Set(ctx, "message:"+strconv.FormatUint(uint64(id), 10), message, 0)

	if err := set.Err(); err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (repo *RedisMessageRepository) Fetch(id uint) (string, bool, error) {
	client := repo.client
	ctx := repo.ctx

	message, err := client.Get(ctx, "message:"+strconv.FormatUint(uint64(id), 10)).Result()

	if err == redis.Nil {
		return "", false, nil
	} else if err != nil {
		return "", false, err
	}

	return message, true, nil
}
