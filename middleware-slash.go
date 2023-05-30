package framework

import "github.com/gofiber/fiber/v2"

func middlewareSlash() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Next()
	}
}
