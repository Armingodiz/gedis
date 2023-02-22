package ttlManager

import (
	"fmt"
	"gedis/config"
	"time"

	"github.com/go-redis/redis"
)

var singeletonManager Manager

func GetManager() (Manager, error) {
	if singeletonManager == nil {
		err := connectRedis()
		return singeletonManager, err
	}
	return singeletonManager, nil
}

type Manager interface {
	SetTtl(userName, key string, ttl int) error
	CheckKey(userName, key string) (bool, error)
}

type redisManager struct {
	redisClient *redis.Client
}

func connectRedis() error {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Configs.Cache.CacheAddr,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return err
	}
	singeletonManager = &redisManager{redisClient: client}
	return nil
}

func (r *redisManager) SetTtl(userName, key string, ttl int) error {
	if ttl == -1 {
		return r.redisClient.Set(fmt.Sprintf("%s_%s", userName, key), 1, 0).Err()
	} else {
		if ok, err := r.CheckKey(userName, key); err == nil && ok {
			return r.redisClient.Expire(fmt.Sprintf("%s_%s", userName, key), time.Second*time.Duration(ttl)).Err()
		}
		return r.redisClient.Set(fmt.Sprintf("%s_%s", userName, key), 1, time.Second*time.Duration(ttl)).Err()
	}
}

func (r *redisManager) CheckKey(userName, key string) (bool, error) {
	status := r.redisClient.Exists(fmt.Sprintf("%s_%s", userName, key))
	return status.Val() != 0, status.Err()
}
