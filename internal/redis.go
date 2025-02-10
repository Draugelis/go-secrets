package internal

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisService defines the interface for Redis operations.
type RedisService interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	NewScanner(ctx context.Context, match string) (RedisScanner, error)
	NewPipeline(ctx context.Context) (RedisPipeline, error)
}

type RedisScanner interface {
	Next(ctx context.Context) bool
	Val() string
	Err() error
}

type RedisScannerImpl struct {
	Iterator *redis.ScanIterator
}

type RedisPipeline interface {
	Del(ctx context.Context, key string) error
	Exec(ctx context.Context) ([]RedisResult, error)
	Discard()
}

type RedisPipelineImpl struct {
	Pipe redis.Pipeliner
}

type RedisResult struct {
	Val interface{}
	Err error
}

// Redis represents a struct that holds the Redis client for performing Redis operations.
type RedisServiceImpl struct {
	Client *redis.Client
}

var redisOnce sync.Once
var redisInstance RedisService

func GetRedisService() (RedisService, error) {
	if redisInstance == nil {
		return nil, fmt.Errorf("RedisService is not initialized, call SetupRedis first")
	}
	return redisInstance, nil
}

func SetupRedis(address string) error {
	var err error
	redisOnce.Do(func() {
		redisInstance, err = NewRedisService(address)
	})
	return err
}

// NewRedisService initializes a Redis client connection and returns it.
func NewRedisService(address string) (*RedisServiceImpl, error) {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisServiceImpl{Client: client}, nil
}

// Set stores a value in Redis with a specified TTL.
func (r *RedisServiceImpl) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := r.Client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("could not set key: %w", err)
	}
	slog.Info("key set in redis", "key", key)
	return nil
}

// Get retrieves a value from Redis.
func (r *RedisServiceImpl) Get(ctx context.Context, key string) (string, error) {
	value, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key does not exist: %s", key)
	}
	if err != nil {
		return "", fmt.Errorf("could not get key: %w", err)
	}
	return value, nil
}

// Del removes a key from Redis.
func (r *RedisServiceImpl) Del(ctx context.Context, key string) error {
	err := r.Client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("could not delete key: %w", err)
	}
	return nil
}

// TTL retrieves the TTL (time to live) for a key in Redis.
func (r *RedisServiceImpl) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.Client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("could not get TTL for key: %w", err)
	}
	return ttl, nil
}

func (r *RedisServiceImpl) NewScanner(ctx context.Context, match string) (RedisScanner, error) {
	iter := r.Client.Scan(ctx, 0, match, 100).Iterator()
	return &RedisScannerImpl{Iterator: iter}, nil
}

func (rs *RedisScannerImpl) Next(ctx context.Context) bool {
	return rs.Iterator.Next(ctx)
}

func (rs *RedisScannerImpl) Val() string {
	return rs.Iterator.Val()
}

func (rs *RedisScannerImpl) Err() error {
	return rs.Iterator.Err()
}

func (r *RedisServiceImpl) NewPipeline(ctx context.Context) (RedisPipeline, error) {
	pipe := r.Client.Pipeline()
	return &RedisPipelineImpl{Pipe: pipe}, nil
}

func (rp *RedisPipelineImpl) Del(ctx context.Context, key string) error {
	return rp.Pipe.Del(ctx, key).Err()
}

func (rp *RedisPipelineImpl) Exec(ctx context.Context) ([]RedisResult, error) {
	cmds, err := rp.Pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("pipeline execution failed: %w", err)
	}

	results := make([]RedisResult, len(cmds))
	for i, cmd := range cmds {
		stringCmd, ok := cmd.(*redis.StringCmd)
		if !ok {
			return nil, fmt.Errorf("unexpected command type: %T", cmd)
		}
		results[i] = RedisResult{Val: stringCmd.Val(), Err: stringCmd.Err()}
	}

	return results, nil
}

func (rp *RedisPipelineImpl) Discard() {
	rp.Pipe.Discard()
}
