package framework

import (
	"github.com/gofiber/fiber/v2"
)

type control struct {
	ctx *fiber.Ctx
}

func newControl(ctx *fiber.Ctx) *control {
	return &control{
		ctx: ctx,
	}
}

func (c *control) Path() string {
	return c.ctx.Route().Path
}
