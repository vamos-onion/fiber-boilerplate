package middleware

import (
	"strings"

	"fiber-boilerplate/internal/app/config"
	"fiber-boilerplate/internal/app/handlers"
	logging "fiber-boilerplate/internal/pkg/logging"
	"fiber-boilerplate/internal/pkg/session"
	"fiber-boilerplate/internal/pkg/setting"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Register(f *fiber.App) {
	f.Use(requestid.New())
	f.Use(logger.New())
	f.Use(recover.New())
	f.Use(etag.New())
	f.Use(compress.New())

	// Only enable pprof in non-production environments
	if setting.Runtime.Env == "local" || setting.Runtime.Env == "development" {
		f.Use(pprof.New())
		logging.Info("Pprof profiling enabled at /debug/pprof (environment: %s)", setting.Runtime.Env)
	}

	// Configure CORS
	if config.Server.CORS.Enabled {
		f.Use(cors.New(cors.Config{
			AllowOrigins:     strings.Join(config.Server.CORS.AllowOrigins, ","),
			AllowMethods:     strings.Join(config.Server.CORS.AllowMethods, ","),
			AllowHeaders:     strings.Join(config.Server.CORS.AllowHeaders, ","),
			ExposeHeaders:    strings.Join(config.Server.CORS.ExposeHeaders, ","),
			MaxAge:           config.Server.CORS.MaxAge,
			AllowCredentials: config.Server.CORS.AllowCredentials,
		}))
	}

	f.Use(oapiRequestValidate)

	f.Use(keyauth.New(keyauth.Config{
		KeyLookup:  "header:Authorization",
		AuthScheme: "Bearer",
		Validator:  validateAPIKey,
		Next: func(c *fiber.Ctx) bool {
			// OpenAPI 검증에서 security가 없는 경로는 JWT 검증 스킵
			// oapiRequestValidate에서 설정한 컨텍스트 값 확인
			skipAuth, ok := c.Locals("oapi:skip_auth").(bool)
			logging.Debug("keyauth Next: path=%s, skipAuth=%v, ok=%v", c.Path(), skipAuth, ok)
			if ok && skipAuth {
				// true를 반환하면 이 미들웨어를 스킵
				logging.Debug("Skipping keyauth for path: %s", c.Path())
				return true
			}
			// false를 반환하면 JWT 검증 수행
			logging.Debug("Running keyauth for path: %s", c.Path())
			return false
		},
	}))

	sessionMiddleware := session.Middleware(ContextKeyStore, validate, handlers.SendError)
	f.Use(func(ctx *fiber.Ctx) error {
		handler := sessionMiddleware(func(c *fiber.Ctx) error {
			return c.Next()
		})
		return handler(ctx)
	})
}
