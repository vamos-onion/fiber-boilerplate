package database

import (
	"encoding/json"

	"fiber-boilerplate/internal/pkg/util"
)

// DriverConfigBlock :
type DriverConfigBlock struct {
	// PostgreSQL fields
	Conn         string `env:"DB_CONN" json:"conn,omitempty"`
	Host         string `env:"DB_HOST" json:"host,omitempty"`
	Port         int    `env:"DB_PORT" envDefault:"5432" json:"port,omitempty"`
	User         string `env:"DB_USER" json:"user,omitempty"`
	Password     string `env:"DB_PASSWORD" envDefault:"" json:"password,omitempty"`
	Name         string `env:"DB_NAME" envDefault:"playground" json:"name,omitempty"`
	MaxIdleConns int    `env:"DB_MAX_IDLE_CONNS" envDefault:"5" json:"maxIdleConns,omitempty"`
	MaxOpenConns int    `env:"DB_MAX_OPEN_CONNS" envDefault:"100" json:"maxOpenConns,omitempty"`
	MaxLifetime  int    `env:"DB_MAX_LIFETIME" envDefault:"300" json:"maxLifeTime,omitempty"`

	// Redis fields
	RedisHost     string `env:"REDIS_HOST" json:"redisHost,omitempty"`
	RedisPort     int    `env:"REDIS_PORT" envDefault:"6379" json:"redisPort,omitempty"`
	RedisPassword string `env:"REDIS_PASSWORD" json:"redisPassword,omitempty"`
}

// DriverEnum : 드라이버 Enum
type DriverEnum int

// DriverEnums
const (
	DriverPostgres DriverEnum = iota // postgres
	DriverRedis                      // redis
	driverCount
)

// DriverEnumNames :
var DriverEnumNames = [driverCount]string{"postgres", "redis"}

func (id DriverEnum) String() string {
	return DriverEnumNames[id]
}

var driverConfigs [driverCount]DriverConfigBlock

// Setup :
func Setup(rawDatabases json.RawMessage) {
	var databases map[string]json.RawMessage
	err := util.Json.Unmarshal(rawDatabases, &databases)
	if err != nil {
		panic(err)
	}

	for id, name := range DriverEnumNames {
		raw := databases[name]
		if raw != nil {
			var config DriverConfigBlock
			err := util.Json.UnmarshalWithEnv(raw, &config)
			if err != nil {
				panic(err)
			}

			driverConfigs[id] = config
		}
	}
}
