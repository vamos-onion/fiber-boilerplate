package config

import (
	"encoding/json"

	"fiber-boilerplate/internal/pkg/database"
	"fiber-boilerplate/internal/pkg/setting"
	"fiber-boilerplate/internal/pkg/util"
)

// ServerBlock : admin ServerBlock
type ServerBlock struct {
	GracefulTimeout    int    `env:"GRACEFUL_TIMEOUT" envDefault:"10" json:"gracefulTimeout,omitempty"`
	LogRequestsEnabled bool   `env:"LOG_REQUESTS_ENABLED" envDefault:"true" json:"logRequestsEnabled,omitempty"`
	JwtSecret          string `env:"JWT_SECRET" json:"jwtSecret,omitempty"`
	CORS               struct {
		Enabled          bool     `env:"CORS_ENABLED" envDefault:"true" json:"enabled,omitempty"`
		AllowOrigins     []string `env:"CORS_ALLOW_ORIGINS" envSeparator:"," envDefault:"*" json:"allowOrigins,omitempty"`
		AllowMethods     []string `env:"CORS_ALLOW_METHODS" envSeparator:"," envDefault:"GET,HEAD,PUT,PATCH,POST,DELETE,OPTIONS" json:"allowMethods,omitempty"`
		AllowHeaders     []string `env:"CORS_ALLOW_HEADERS" envSeparator:"," envDefault:"Cache-Control" json:"allowHeaders,omitempty"`
		AllowCredentials bool     `env:"CORS_ALLOW_CREDENTIALS" envDefault:"false" json:"allowCredentials,omitempty"`
		ExposeHeaders    []string `env:"CORS_EXPOSE_HEADERS" envSeparator:"," envDefault:"" json:"exposeHeaders,omitempty"`
		MaxAge           int      `env:"CORS_MAX_AGE" envSeparator:"," envDefault:"0" json:"maxAge,omitempty"`
	} `json:"cors"`
}

// Server : admin server
var Server ServerBlock

// Setup : admin setup
func Setup() {
	err := util.Json.UnmarshalWithEnv(nil, &Server)
	if err != nil {
		panic(err)
	}

	// Validate required configuration
	if Server.JwtSecret == "" {
		panic("JWT_SECRET environment variable is required")
	}

	// Setup database configuration
	// Both postgres and redis fields will be read from environment variables
	databaseConfig := json.RawMessage(`{ "postgres": {}, "redis": {} }`)
	if setting.Configs == nil {
		setting.Configs = make(map[string]json.RawMessage)
	}
	setting.Configs["databases"] = databaseConfig
	database.Setup(setting.Configs["databases"])
}
