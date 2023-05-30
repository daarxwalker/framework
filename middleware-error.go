package framework

import (
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/mailgun/raymond/v2"
)

func middlewareError(app *App) fiber.Handler {
	return func(ctx *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				recoverType := reflect.TypeOf(r)
				if recoverType != errorType {
					return
				}
				e := r.(managedError)
				if e.status == 0 {
					e.status = fiber.StatusBadRequest
				}

				// Template
				errorTmplName := fmt.Sprintf("%d", e.status)
				tmpl, ok := app.templatesManager.errors[errorTmplName]
				if !ok {
					err = ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("template [%s] not found", errorTmplName))
					return
				}
				html, err := tmpl.Exec(e)
				if err != nil {
					err = ctx.Status(fiber.StatusInternalServerError).SendString(e.message)
					return
				}

				// Layout
				layoutTmpl, ok := app.templatesManager.layouts[e.control.layout]
				if ok {
					contextData := make(map[string]any)
					contextData["router"] = func(options *raymond.Options) raymond.SafeString {
						routeHtml := html
						return raymond.SafeString(routeHtml)
					}
					html, err = layoutTmpl.Exec(contextData)
				}

				ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
				err = ctx.Status(fiber.StatusOK).SendString(html)
			}
		}()
		return ctx.Next()
	}
}
