package session

import (
	"fiber-boilerplate/internal/defs"
	"fiber-boilerplate/internal/models"
	logging "fiber-boilerplate/internal/pkg/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ContextKeyStore = "fiber-boilerplate/#/sessionStore"
	ContextKeyData  = "fiber-boilerplate/#/sessionData"
)

// DataBlock :
type DataBlock struct {
	Appuser *models.AppuserBlock
}

// New :
func New(args ...interface{}) *DataBlock {
	data := new(DataBlock)

	for _, arg := range args {
		switch arg.(type) {
		case *models.AppuserBlock:
			data.Appuser = arg.(*models.AppuserBlock)

		case *models.EntityBlock:
			e := arg.(*models.EntityBlock)
			if e.Appuser != nil {
				data.Appuser = e.Appuser
			}

		default:
			panic(defs.ErrInvalid)
		}
	}

	return data
}

// Create :
func Create(ctx *fiber.Ctx, key string, data *DataBlock) error {
	store := ctx.Locals(ContextKeyStore).(*StoreBlock)

	err := store.Set(ctx.Context(), key, data)
	if err == nil {
		logging.Trace("session created: %s %+v", key, data)
	}

	return err
}

// Update :
func Update(ctx *fiber.Ctx, data *DataBlock) error {
	store := ctx.Locals(ContextKeyStore).(*StoreBlock)
	key := ctx.Locals(store.KeyName).(jwt.MapClaims)["uuid"].(string)

	err := store.Set(ctx.Context(), key, data)
	if err == nil {
		logging.Trace("session updated: %s %+v", key, data)
	}

	return err
}

// FromContext :
func FromContext(ctx *fiber.Ctx) *DataBlock {
	if _, ok := ctx.Locals(ContextKeyData).(*DataBlock); !ok {
		return nil
	}
	return ctx.Locals(ContextKeyData).(*DataBlock)
}

// Destroy :
func Destroy(ctx *fiber.Ctx) error {
	store := ctx.Locals(ContextKeyStore).(*StoreBlock)
	key := ctx.Locals(store.KeyName).(jwt.MapClaims)["uuid"].(string)

	err := store.Del(ctx.Context(), key)
	if err == nil {
		logging.Trace("session destroyed: %s", key)
	}

	return store.Del(ctx.Context(), key)
}

// GetStore :
func GetStore(keyName string) *StoreBlock {
	return newStore(keyName)
}
