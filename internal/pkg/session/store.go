package session

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"time"

	"fiber-boilerplate/internal/defs"
	"fiber-boilerplate/internal/pkg/database"
	logging "fiber-boilerplate/internal/pkg/logging"
)

// StoreBlock :
type StoreBlock struct {
	KeyName string
	redis   *database.Redis
}

// CloseConnection : Renamed from Close to prevent automatic cleanup by Fiber
func (s *StoreBlock) CloseConnection() error {
	if s == nil || s.redis == nil {
		return defs.ErrInvalid
	}
	return s.redis.Close()
}

// Del :
func (s *StoreBlock) Del(ctx context.Context, key string) error {
	if s == nil || s.redis == nil {
		return defs.ErrInvalid
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return s.redis.Del(ctx, key)
}

// Set :
func (s *StoreBlock) Set(ctx context.Context, key string, data *DataBlock) error {
	if s == nil || s.redis == nil {
		return defs.ErrInvalid
	}
	if ctx == nil {
		ctx = context.Background()
	}

	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(data)
	if err != nil {
		return err
	}

	return s.redis.Set(ctx, key, buf.Bytes())
}

// Get :
func (s *StoreBlock) Get(ctx context.Context, key string) (*DataBlock, bool, error) {
	if s == nil || s.redis == nil {
		return nil, false, defs.ErrInvalid
	}
	if ctx == nil {
		ctx = context.Background()
	}

	value, found, err := s.redis.Get(ctx, key)
	if err != nil {
		return nil, false, err

	} else if found {
		var data DataBlock
		var buf bytes.Buffer

		buf.Write(value.([]byte))
		decoder := gob.NewDecoder(&buf)
		err = decoder.Decode(&data)
		if err != nil {
			logging.Warn(err, "Failed to decode session data for key: %s. Deleting corrupted session.", key)
			_ = s.redis.Del(ctx, key) // Clean up the invalid session
			return nil, false, nil
		}

		return &data, true, nil
	}

	return nil, false, nil
}

// Flush :
func (s *StoreBlock) Flush(ctx context.Context) error {
	if s == nil || s.redis == nil {
		return defs.ErrInvalid
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return s.redis.Flush(ctx)
}

// Close :
func (s *StoreBlock) Close() error {
	if s == nil || s.redis == nil {
		return defs.ErrInvalid
	}
	return s.redis.Close()
}

// SubscribeChannel :
func (s *StoreBlock) SubscribeChannel(ctx context.Context, channel string) error {
	if s == nil || s.redis == nil {
		return defs.ErrInvalid
	}

	if ctx == nil {
		ctx = context.Background()
	}

	ctxWTO, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Subscriber Goroutine
	go func() {
		sub := s.redis.SubscribeChannel(ctx, channel)
		ch := sub.Channel()
		defer sub.Close()

		for {
			select {
			case <-ctxWTO.Done():
				// Context canceled or timeout, exit goroutine
				logging.Trace("Context canceled or timeout, exit goroutine")
				return
			case msg := <-ch:
				if msg == nil {
					// Channel closed
					logging.Trace("Channel closed")
					return
				}
				logging.Trace("Received message:", msg.Payload)
				// 메시지 받으면 종료
				if len(ch) == 0 {
					logging.Trace("Received all messages, exit goroutine")
					// context cancel
					cancel()
					return
				}
			}
		}
	}()

	// Main Goroutine waits
	select {
	case <-ctxWTO.Done():
		// Timeout or cancellation
		if errors.Is(ctxWTO.Err(), context.DeadlineExceeded) {
			logging.Trace("Context timeout")
			return defs.ErrTimeout
		}

		logging.Trace("Context Done")
		return nil
	}
}

// PublishChannel :
func (s *StoreBlock) PublishChannel(ctx context.Context, channel string, message string) error {
	if s == nil || s.redis == nil {
		return defs.ErrInvalid
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return s.redis.PublishChannel(ctx, channel, message)
}

func newStore(keyName string) *StoreBlock {
	var db int
	var ttlSec int

	switch keyName {
	case "fiber-boilerplate/#/appuserID", "fiber-boilerplate/#/keyStore":
		db = 0
		ttlSec = 7200 // 2시간
	default:
		panic(defs.ErrInvalid)
	}

	logging.Info("Creating new store for keyName: %s", keyName)

	return &StoreBlock{
		KeyName: keyName,
		redis:   database.NewRedis(db, ttlSec),
	}
}
