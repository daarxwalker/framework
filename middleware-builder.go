package framework

import (
	"strings"

	"dd"
)

type MiddlewareBuilder struct {
	app      *App
	ctx      *Ctx
	name     string
	handler  func() error
	path     string
	langCode string
}

func CreateMiddleware(name string, app *App, ctx *Ctx) *MiddlewareBuilder {
	b := &MiddlewareBuilder{
		app:  app,
		name: name,
		ctx:  ctx,
	}
	b.prepare()
	dd.Print(b)
	return b
}

func (b *MiddlewareBuilder) LangCode() string {
	if !b.app.i18n {
		return ""
	}
	return b.langCode
}

func (b *MiddlewareBuilder) Language() (Language, bool) {
	l, ok := b.app.languages[b.langCode]
	return l, ok
}

func (b *MiddlewareBuilder) MainLanguage() (Language, bool) {
	for _, l := range b.app.languages {
		if l.main {
			return l, true
		}
	}
	return &language{}, false
}

func (b *MiddlewareBuilder) LangExist() bool {
	return len(b.langCode) > 0
}

func (b *MiddlewareBuilder) IsLangValid() bool {
	var valid bool
	for code := range b.app.languages {
		if strings.HasPrefix(b.path, "/"+code+"/") {
			valid = true
		}
	}
	return valid
}

func (b *MiddlewareBuilder) Redirect(path string) error {
	return b.ctx.Redirect(path)
}

func (b *MiddlewareBuilder) Next() error {
	return b.ctx.Next()
}

func (b *MiddlewareBuilder) prepare() {
	b.path = b.ctx.Path()
	b.prepareLangCode()
}

func (b *MiddlewareBuilder) prepareLangCode() {
	if !b.app.i18n {
		return
	}
	b.langCode = b.ctx.Params(langParam)
	if len(b.langCode) > 0 || len(b.path) < 4 {
		return
	}
	b.langCode = b.path[1:3]
}
