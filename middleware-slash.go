package framework

import (
	"strings"
)

func middlewareSlash(app *App) MiddlewareHandler {
	return func(ctx *Ctx) error {
		m := CreateMiddleware("slash", app, ctx)
		if !strings.HasSuffix(m.path, "/") {
			return m.Redirect(m.path + "/")
		}
		return m.Next()
	}
}
