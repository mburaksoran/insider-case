package repository

import (
	"context"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/mburaksoran/insider-case/internal/domain/models"
	"github.com/mburaksoran/insider-case/internal/domain/repository"
	"github.com/mburaksoran/insider-case/internal/infra/engines"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
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
		return errors.New("[Set] - Could not obtain lock!")
	}
	defer lock.Release(ctx)
	err := r.RedisClient.Set(ctx, key, value, TTL).Err()
	if err != nil {
		r.logger.Errorf("[Set] - Error on setting messageReceiveHistory Err: %s", err.Error())
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
	if errors.Is(err, redislock.ErrNotObtained) {
		r.logger.Warnf("[lockRedis] - Could not obtain lock!")
		return nil
	} else if err != nil {
		r.logger.Errorf("[lockRedis] - Error on obtaining lock Err: %s", err.Error())
		return nil
	}
	return lock
}
