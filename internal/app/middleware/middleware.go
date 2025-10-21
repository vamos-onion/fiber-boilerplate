package middleware

import (
	"fmt"
	"runtime/debug"
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
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
)

// LoggerWriter : Fiber logger의 출력을 zerolog로 전달하는 커스텀 writer
// Custom io.Writer that routes Fiber logger output to zerolog Info level
type LoggerWriter struct{}

func (lw *LoggerWriter) Write(p []byte) (n int, err error) {
	logging.Info(string(p))
	return len(p), nil
}

func Register(f *fiber.App) {
	// Request ID: UUID로 각 요청 추적 (X-Request-ID 헤더)
	// Generates unique IDs for request tracking and distributed tracing
	f.Use(requestid.New(requestid.Config{
		Header: "X-Request-ID",
		Generator: func() string {
			return uuid.New().String()
		},
	}))

	// Logger: 모든 HTTP 요청/응답을 zerolog로 로깅 (접근 로그)
	// Logs all incoming HTTP requests via zerolog with detailed information
	f.Use(logger.New(logger.Config{
		Format:   "${time} | ${status} | ${method} | ${path} | ${latency}ms | ClientIP: ${ip} - RequestID: ${locals:requestid}\n",
		TimeZone: "Local",
		Output:   &LoggerWriter{}, // LoggerWriter를 통해 zerolog로 라우팅
	}))

	// Recover: 패닉 발생 시 자동 복구 (스택 트레이스 로깅, JSON 에러 응답)
	// Catches panics with stack trace logging and returns JSON error response
	f.Use(func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				panicErr := fmt.Errorf("internal server error: %v", r)
				logging.Error(panicErr, "Panic: %v\nPath: %s\nMethod: %s\nStack:\n%s",
					fmt.Sprintf("%v", r), c.Path(), c.Method(), debug.Stack())
				handlers.SendError(c, fiber.StatusInternalServerError, panicErr)
			}
		}()
		return c.Next()
	})

	// ETag: Strong ETag 생성으로 HTTP 캐싱 지원 (304 Not Modified 활용)
	// Generates ETags for efficient HTTP caching and bandwidth reduction
	f.Use(etag.New(etag.Config{
		Weak: false,
	}))

	// Compress: gzip/deflate/brotli로 응답 압축 (기본 압축 레벨)
	// Compresses responses to reduce bandwidth usage
	f.Use(compress.New(compress.Config{
		Level: compress.LevelDefault,
	}))

	// Pprof: 성능 프로파일링 (개발/로컬 환경만, /debug/pprof)
	// Performance profiling endpoint (local/development only)
	if setting.Runtime.Env == "local" || setting.Runtime.Env == "development" {
		f.Use(pprof.New())
		logging.Info("Pprof profiling enabled at /debug/pprof (environment: %s)", setting.Runtime.Env)
	}

	// CORS: 크로스 오리진 요청 처리 (환경변수로 설정)
	// Handles cross-origin requests based on configuration
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

	// OpenAPI Validation: 요청이 OpenAPI 스펙을 준수하는지 검증
	// Validates requests against OpenAPI specification (body, params, headers)
	f.Use(oapiRequestValidate)

	// Key Auth (JWT): Authorization 헤더에서 Bearer 토큰 추출 및 검증
	// Extracts and validates JWT tokens from Authorization header
	// OpenAPI의 security 필드에 따라 인증 스킵 여부 결정 (oapi:skip_auth)
	f.Use(keyauth.New(keyauth.Config{
		KeyLookup:  "header:Authorization",
		AuthScheme: "Bearer",
		Validator:  validateAPIKey,
		Next: func(c *fiber.Ctx) bool {
			skipAuth, ok := c.Locals("oapi:skip_auth").(bool)
			logging.Debug("keyauth Next: path=%s, skipAuth=%v, ok=%v", c.Path(), skipAuth, ok)
			if ok && skipAuth {
				logging.Debug("Skipping keyauth for path: %s", c.Path())
				return true
			}
			logging.Debug("Running keyauth for path: %s", c.Path())
			return false
		},
	}))

	// Session: JWT의 uuid 클레임으로 DB에서 사용자 정보 로드
	// Loads user session from database using JWT uuid claim
	sessionMiddleware := session.Middleware(ContextKeyStore, validate, handlers.SendError)
	f.Use(func(ctx *fiber.Ctx) error {
		handler := sessionMiddleware(func(c *fiber.Ctx) error {
			return c.Next()
		})
		return handler(ctx)
	})
}
