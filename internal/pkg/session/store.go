package session

import (
	"bytes"
	"context"
	"encoding/gob"

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
	return s.redis.Close()
}

// Del :
func (s *StoreBlock) Del(ctx context.Context, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	return s.redis.Del(ctx, key)
}

// Set :
func (s *StoreBlock) Set(ctx context.Context, key string, data *DataBlock) error {
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
	if ctx == nil {
		ctx = context.Background()
	}
	return s.redis.Flush(ctx)
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
