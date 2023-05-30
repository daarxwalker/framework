package framework

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func middlewareI18nKnownLangs(app *App) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if !app.i18n {
			return ctx.Next()
		}
		lang := ctx.Params(langParam)
		if _, ok := app.languages[lang]; !ok {
			for _, l := range app.languages {
				if l.main {
					return ctx.Redirect(strings.Replace(ctx.Path(), lang, l.code, -1))
				}
			}
		}
		return ctx.Next()
	}
}

func middlewareI18nAddLangUrlPrefix(app *App) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if !app.i18n {
			return ctx.Next()
		}
		var valid bool
		var mainLang *language
		for code, l := range app.languages {
			if l.main {
				mainLang = l
			}
			if strings.HasPrefix(ctx.Path(), "/"+code+"/") {
				valid = true
			}
		}
		if valid {
			return ctx.Next()
		}
		if mainLang == nil {
			return ctx.Next()
		}
		path := ctx.Path()
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		path = "/" + mainLang.code + path
		return ctx.Redirect(path)
	}
}
