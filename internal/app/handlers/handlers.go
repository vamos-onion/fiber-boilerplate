package handlers

import (
	v1 "fiber-boilerplate/internal/app/handlers/v1"
	api "fiber-boilerplate/internal/generated/serviceapi"

	"github.com/gofiber/fiber/v2"
)

// APIHandlerBlock :
type APIHandlerBlock struct{}

// SendError :
func SendError(ctx *fiber.Ctx, code int, errs ...error) error {
	return v1.SendError(ctx, code, errs...)
}

// SendResponse :
func SendResponse(ctx *fiber.Ctx, code int, response interface{}) error {
	return v1.SendResponse(ctx, code, response)
}

func (h APIHandlerBlock) GetPing(ctx *fiber.Ctx) error {
	return v1.GetPing(ctx)
}

func (h APIHandlerBlock) SseOpen(ctx *fiber.Ctx) error {
	return v1.SseOpen(ctx)
}

func (h APIHandlerBlock) SseClose(ctx *fiber.Ctx) error {
	return v1.SseClose(ctx)
}

func (h APIHandlerBlock) ListAppusers(ctx *fiber.Ctx, params api.ListAppusersParams) error {
	return v1.ListAppusers(ctx, params)
}

func (h APIHandlerBlock) CreateAppuser(ctx *fiber.Ctx) error {
	return v1.CreateAppuser(ctx)
}

func (h APIHandlerBlock) UpdateAppuser(ctx *fiber.Ctx) error {
	return v1.UpdateAppuser(ctx)
}
