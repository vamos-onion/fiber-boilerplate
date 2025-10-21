/*
	Redis 6.x는 v8 패키지 사용
	Redis 7.x는 v9 패키지 사용
	GCP Redis 는 6.x 이므로 v8 패키지를 사용하도록 한다
	https://github.com/go-redis/redis#installation

	Distributed Lock Manager 는 아래 참고
	https://redis.io/docs/reference/patterns/distributed-locks/
*/

package database

import (
	"context"
	"time"

	"fiber-boilerplate/internal/pkg/cache"
	logging "fiber-boilerplate/internal/pkg/logging"
	"fiber-boilerplate/internal/pkg/setting"
	"fiber-boilerplate/internal/pkg/util"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

// Redis :
type Redis struct {
	Set              func(ctx context.Context, key string, value interface{}) error
	Get              func(ctx context.Context, key string) (interface{}, bool, error)
	Del              func(ctx context.Context, key string) error
	Flush            func(ctx context.Context) error
	SubscribeChannel func(ctx context.Context, channel string) *redis.PubSub
	PublishChannel   func(ctx context.Context, channel string, message string) error
	Close            func() error
}

func sharedKey(key string) string {
	return util.String.Concat("shared/", key)
}

// NewRedis :
func NewRedis(db int, ttlSec int) *Redis {
	config := driverConfigs[DriverRedis]
	conn := config.Conn

	// Build connection string from RedisHost and RedisPort if Conn is empty
	if len(conn) < 1 {
		if len(config.RedisHost) > 0 {
			conn = util.String.Concat(config.RedisHost, ":", config.RedisPort)
		}
	}

	r := new(Redis)

	ttl := time.Duration(ttlSec) * time.Second

	cli := redis.NewClient(&redis.Options{
		Addr:         conn,
		Password:     config.RedisPassword,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     20,
		MinIdleConns: 10,
		MaxRetries:   3,
		PoolTimeout:  4 * time.Second,
	})

	// GCP 환경에서의 redis status check
	var err error
	if setting.Runtime.Env == "development" || setting.Runtime.Env == "staging" || setting.Runtime.Env == "production" {
		// timeout
		c, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()
	LOOP:
		for trial := 0; trial < 5; trial++ {
			err = cli.Ping(c).Err()
			switch err {
			case nil:
				break LOOP
			default:
				logging.Warn(err, "redis trial:%d", trial)
				time.Sleep(time.Second)
			}
			if trial == 4 {
				panic(err)
			}
		}
	} else {
		err = cli.Ping(context.Background()).Err()
	}

	if err == nil {
		logging.Info("Redis: connected to %s, database %d, ttl %d", conn, db, ttlSec)

		pool := goredis.NewPool(cli)
		sync := redsync.New(pool)
		mutex := sync.NewMutex("fiber-boilerplate-redis-lock")

		r.Set = func(ctx context.Context, key string, value interface{}) error {
			if ctx == nil {
				ctx = context.Background()
			}
			err := mutex.LockContext(ctx)
			if err != nil {
				return err
			}
			defer func() {
				_, _ = mutex.UnlockContext(ctx)
			}()

			sKey := sharedKey(key)
			return cli.Set(ctx, sKey, value, ttl).Err()
		}

		r.Get = func(ctx context.Context, key string) (interface{}, bool, error) {
			if ctx == nil {
				ctx = context.Background()
			}
			err := mutex.LockContext(ctx)
			if err != nil {
				return nil, false, err
			}
			defer func() {
				_, _ = mutex.UnlockContext(ctx)
			}()

			sKey := sharedKey(key)
			value, err := cli.Get(ctx, sKey).Bytes()
			switch err {
			case redis.Nil:
				return nil, false, nil
			case nil:
				return value, true, nil
			default:
				return nil, false, err
			}
		}

		r.Del = func(ctx context.Context, key string) error {
			if ctx == nil {
				ctx = context.Background()
			}
			err := mutex.LockContext(ctx)
			if err != nil {
				return err
			}
			defer func() {
				_, _ = mutex.UnlockContext(ctx)
			}()

			sKey := sharedKey(key)
			return cli.Del(ctx, sKey).Err()
		}

		r.Flush = func(ctx context.Context) error {
			if ctx == nil {
				ctx = context.Background()
			}
			err := mutex.LockContext(ctx)
			if err != nil {
				return err
			}
			defer func() {
				_, _ = mutex.UnlockContext(ctx)
			}()

			return cli.FlushDB(ctx).Err()
		}

		r.SubscribeChannel = func(ctx context.Context, channel string) *redis.PubSub {
			if ctx == nil {
				ctx = context.Background()
			}
			err := mutex.LockContext(ctx)
			if err != nil {
				return nil
			}
			defer func() {
				_, _ = mutex.UnlockContext(ctx)
			}()

			return cli.Subscribe(ctx, channel)
		}

		r.PublishChannel = func(ctx context.Context, channel string, message string) error {
			if ctx == nil {
				ctx = context.Background()
			}
			err := mutex.LockContext(ctx)
			if err != nil {
				return err
			}
			defer func() {
				_, _ = mutex.UnlockContext(ctx)
			}()

			return cli.Publish(ctx, channel, message).Err()
		}

		r.Close = func() error {
			logging.Info("Redis: closing connection, database %d (this should not happen during normal operation)", db)
			return cli.Close()
		}

	} else {
		logging.Warn(err, "Redis: fallback to in-memory cache")

		fallback := cache.New(cache.StoreMemoryDefault, ttlSec)

		r.Set = func(ctx context.Context, key string, value interface{}) error {
			return fallback.Set(key, value)
		}
		r.Get = func(ctx context.Context, key string) (interface{}, bool, error) {
			return fallback.Get(key)
		}
		r.Del = func(ctx context.Context, key string) error {
			return fallback.Del(key)
		}
		r.Flush = func(ctx context.Context) error {
			return fallback.Clear()
		}
		// do nothing
		r.SubscribeChannel = nil
		r.PublishChannel = nil
		r.Close = func() error { return nil }
	}

	return r
}
