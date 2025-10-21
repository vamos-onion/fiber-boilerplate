package session

import (
	"net/http"

	logging "fiber-boilerplate/internal/pkg/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// ManagerFunc :
type ManagerFunc func(*fiber.Ctx, string) (int, *DataBlock, error)

// CheckerFunc :
type CheckerFunc func(*fiber.Ctx) error

// MiddlewareFunc :
type MiddlewareFunc func(next fiber.Handler) fiber.Handler

// ErrorHandlerFunc :
type ErrorHandlerFunc func(*fiber.Ctx, int, ...error) error

func Middleware(keyName string, manager ManagerFunc, errorHandler ErrorHandlerFunc) MiddlewareFunc {
	store := newStore(keyName)

	return func(next fiber.Handler) fiber.Handler {
		return func(ctx *fiber.Ctx) error {
			ctx.Locals(ContextKeyStore, store)

			if ctx.Locals(store.KeyName) == nil {
				return next(ctx)
			}

			key, ok := ctx.Locals(store.KeyName).(jwt.MapClaims)
			if !ok {
				logging.Warn(nil, "Invalid session key type: %T. url path: %s", ctx.Locals(store.KeyName), ctx.Request().URI().Path())
				return next(ctx)
			}

			ctx.Locals(store.KeyName, key)
			if key == nil {
				logging.Trace("key is null. url path: %s", ctx.Request().URI().Path())
				// OpenAPI request validator 에서 검사함
				return next(ctx)
			}

			uuidVal, ok := key["uuid"]
			if !ok {
				logging.Warn(nil, "Missing uuid in JWT claims. url path: %s", ctx.Request().URI().Path())
				return errorHandler(ctx, http.StatusUnauthorized)
			}

			uuid, ok := uuidVal.(string)
			if !ok {
				logging.Warn(nil, "Invalid uuid type in JWT claims: %T. url path: %s", uuidVal, ctx.Request().URI().Path())
				return errorHandler(ctx, http.StatusUnauthorized)
			}

			data, found, err := store.Get(ctx.Context(), uuid)
			if err != nil {
				logging.Warn(err, "get key failed. url path: %s", ctx.Request().URI().Path())
				return errorHandler(ctx, http.StatusInternalServerError, err)

			} else if found {
				// 저장돼 있는 세션이 있으면 통과
				ctx.Locals(ContextKeyData, data)

			} else {
				// 저장돼 있는 세션이 없으면 검증
				logging.Trace("session not found: %s %+v. url path: %s", key, data, ctx.Request().URI().Path())
				// appUserAssert 내에 version checker가 있음
				code, data, err := manager(ctx, uuid)
				if code == http.StatusAccepted {
					// validation api인 경우이고, 세션이 없을 때는 통과
					return next(ctx)
				} else if code != http.StatusOK {
					return errorHandler(ctx, code, err)
				}

				err = Create(ctx, uuid, data)
				if err != nil {
					logging.Trace("session created error. url path: %s", ctx.Request().URI().Path())
					return errorHandler(ctx, http.StatusInternalServerError, err)
				}

				ctx.Locals(ContextKeyData, data)
			}

			return next(ctx)
		}
	}
}
