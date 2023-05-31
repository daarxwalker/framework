package framework

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func middlewareI18nVerifyLang(app *App) MiddlewareHandler {
	return func(ctx *Ctx) error {
		m := CreateMiddleware("verify-lang", app, ctx)
		if !app.i18n {
			return m.Next()
		}
		if _, ok := m.Language(); ok {
			return m.Next()
		}
		l, ok := m.MainLanguage()
		if !ok {
			return m.Next()
		}
		return m.Redirect(strings.Replace(m.path, m.LangCode(), l.Code(), -1))
	}
}

func middlewareI18nMissingLang(app *App) MiddlewareHandler {
	return func(ctx *fiber.Ctx) error {
		m := CreateMiddleware("missing-lang", app, ctx)
		if !app.i18n {
			return m.Next()
		}
		ml, ok := m.MainLanguage()
		if !ok {
			return m.Next()
		}
		if !m.LangExist() {
			return m.Redirect("/" + ml.Code() + m.path)
		}
		if !m.IsLangValid() {
			return m.Redirect(strings.Replace(m.path, "/"+m.LangCode()+"/", "/"+ml.Code()+"/", -1))
		}
		return m.Next()
	}
}
