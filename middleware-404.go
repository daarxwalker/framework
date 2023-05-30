package framework

import "github.com/gofiber/fiber/v2"

func middleware404() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).SendString("not found")
	}
}
