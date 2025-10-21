package util

import (
	"encoding/json"

	"github.com/caarlos0/env/v6"

	"fiber-boilerplate/internal/defs"
)

type jsonBlock struct{}

// Unmarshal :
func (j jsonBlock) Unmarshal(data interface{}, v interface{}) error {
	switch data.(type) {
	case json.RawMessage:
		return json.Unmarshal(data.(json.RawMessage), v)
	case []byte:
		return json.Unmarshal(data.([]byte), v)
	case string:
		return json.Unmarshal([]byte(data.(string)), v)
	default:
		panic(defs.ErrInvalid)
	}
}

// UnmarshalWithEnv :
func (j jsonBlock) UnmarshalWithEnv(bytes []byte, v interface{}) error {
	err := env.Parse(v)
	if err != nil {
		return err
	}

	if len(bytes) > 0 {
		return j.Unmarshal(bytes, v)
	}

	return nil
}
