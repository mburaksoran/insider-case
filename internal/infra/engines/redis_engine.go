package engines

import (
	"github.com/mburaksoran/insider-case/internal/app/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisEngine struct {
	Client *redis.Client
}

var redisEngine *RedisEngine = nil

func GetRedisEngine() *RedisEngine {
	return redisEngine
}

func SetRedisEngine(cfg *config.AppConfig, lgr *zap.SugaredLogger) error {
	if redisEngine == nil {
		lgr.Info("Setting Redis Engine")
		redisEngine = new(RedisEngine)
		redisEngine.Client = redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
			Password: cfg.Redis.Password,
			DB:       0,
		})
		return nil
	}
	return nil
}
