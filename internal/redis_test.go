package internal

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisService_Get(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	redisService := &Redis{
		Client: db,
	}

	key := "test_key"

	// Success case
	mock.ExpectGet(key).
		SetVal("test_value")

	value, err := redisService.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, "test_value", value)

	// Key not found case
	mock.ExpectGet(key).
		RedisNil()

	value, err = redisService.Get(ctx, key)
	assert.Empty(t, value)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key does not exist")

	// Redis error case
	mock.ExpectGet(key).
		SetErr(fmt.Errorf("redis connection failed"))

	value, err = redisService.Get(ctx, key)
	assert.Empty(t, value)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not get key")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisService_Delete(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	redisService := &Redis{
		Client: db,
	}

	key := "test_key"

	// Success case
	mock.ExpectDel(key).
		SetVal(int64(1))

	err := redisService.Delete(ctx, key)
	assert.NoError(t, err)

	// Redis error case
	mock.ExpectDel(key).
		SetErr(fmt.Errorf("delete failed"))

	err = redisService.Delete(ctx, key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not delete key")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisService_TTL(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	redisService := &Redis{
		Client: db,
	}

	key := "test_key"

	// Success case
	mock.ExpectTTL(key).
		SetVal(time.Second * 3600)

	ttl, err := redisService.TTL(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, time.Second*3600, ttl)

	// Key not found case
	mock.ExpectTTL(key).
		RedisNil()

	ttl, err = redisService.TTL(ctx, key)
	assert.Zero(t, ttl)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not get TTL for key")

	// Redis error case
	mock.ExpectTTL(key).
		SetErr(fmt.Errorf("ttl check failed"))

	ttl, err = redisService.TTL(ctx, key)
	assert.Zero(t, ttl)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not get TTL for key")

	assert.NoError(t, mock.ExpectationsWereMet())
}
