package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/mburaksoran/insider-case/internal/domain/models"
	"github.com/mburaksoran/insider-case/internal/domain/repository"
	"github.com/mburaksoran/insider-case/internal/infra/engines"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"log"
	"time"
)

type redisMessageRepository struct {
	RedisClient *redis.Client
	logger      *zap.SugaredLogger
}

const (
	TTL = time.Duration(24 * time.Hour)
)

func NewRedisMessageRepository(lgr *zap.SugaredLogger) repository.RedisMessageRepositoryInterface {
	redisEngine := engines.GetRedisEngine()
	return &redisMessageRepository{
		RedisClient: redisEngine.Client,
		logger:      lgr,
	}
}

func (r *redisMessageRepository) Set(ctx context.Context, key string, value models.MessageReceiveHistory) error {

	lock := r.lockRedis(ctx)
	if lock == nil {
		fmt.Println(lock)
		return errors.New("Could not obtain lock!")
	}
	defer lock.Release(ctx)
	err := r.RedisClient.Set(ctx, key, value, TTL).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (r *redisMessageRepository) lockRedis(ctx context.Context) *redislock.Lock {
	locker := redislock.New(r.RedisClient)
	backoff := redislock.LimitRetry(redislock.ExponentialBackoff(time.Millisecond*100, time.Millisecond*150), 3)
	lock, err := locker.Obtain(ctx, "lock-key", time.Millisecond*100, &redislock.Options{
		RetryStrategy: backoff,
	})
	if err == redislock.ErrNotObtained {
		fmt.Println("Could not obtain lock!")
		return nil
	} else if err != nil {
		log.Fatalln(err)
		return nil
	}
	return lock
}
