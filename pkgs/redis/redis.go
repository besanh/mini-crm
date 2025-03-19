package redis

import (
	"context"
	"time"

	"github.com/besanh/mini-crm/common/log"
	"github.com/redis/go-redis/v9"
)

type (
	IRedis interface {
		GetClient() *redis.Client
		Connect() error
	}

	RedisClient struct {
		Client *redis.Client
		Config RedisConfig
	}

	RedisConfig struct {
		Host         string
		Password     string
		DB           int
		PoolSize     int
		PoolTimeout  int
		ReadTimeout  int
		WriteTimeout int
	}
)

var Redis IRedis

func NewRedis(config RedisConfig) (IRedis, error) {
	redisClient := &RedisClient{
		Config: config,
	}

	if err := redisClient.Connect(); err != nil {
		log.Error(err)
		return nil, err
	}
	return redisClient, nil
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.Client
}

func (r *RedisClient) Connect() error {
	client := redis.NewClient(&redis.Options{
		Addr:         r.Config.Host,
		Password:     r.Config.Password,
		DB:           r.Config.DB,
		PoolSize:     r.Config.PoolSize,
		PoolTimeout:  time.Duration(r.Config.PoolTimeout) * time.Second,
		ReadTimeout:  time.Duration(r.Config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(r.Config.WriteTimeout) * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return err
	}
	str, err := client.Ping(ctx).Result()
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info(str)
	r.Client = client
	return nil
}
