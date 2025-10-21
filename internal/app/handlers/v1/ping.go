package v1

import "github.com/gofiber/fiber/v2"

func GetPing(ctx *fiber.Ctx) error {
	return SendResponse(ctx, 200, map[string]string{"ping": "pong"})
}
