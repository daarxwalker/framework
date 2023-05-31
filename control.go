package framework

import (
	"github.com/gofiber/fiber/v2"
)

type BaseControl interface {
	LangCode() string
	Languages() []Language
	Root() string
}

type GuardControl interface {
	BaseControl
}

type ComponentControl interface {
	BaseControl
	Form(form *Form)
	Redirect(path string)
}

type ControllerControl interface {
	BaseControl
	JSON(value any)
	Meta() Meta
	Redirect(path string)
	Template(template string)
	Text(value string)
}

type control struct {
	app        *App
	ctx        *fiber.Ctx
	response   *response
	controller *appController
	meta       *metaManager
	module     string
	layout     string
}

func newControl(app *App, ctx *fiber.Ctx, controller *appController, module string) *control {
	c := &control{
		app:        app,
		ctx:        ctx,
		response:   new(response),
		meta:       new(metaManager),
		controller: controller,
		module:     module,
		layout:     "default",
	}
	return c
}

func (c *control) LangCode() string {
	if !c.app.i18n {
		return ""
	}
	langCode := c.ctx.Params(langParam)
	langCodeLen := len(langCode)
	if langCodeLen != 0 {
		return langCode
	}
	if langCodeLen == 0 && len(c.Path()) > 3 {
		langCode = c.ctx.Path()[1:3]
	}
	return langCode
}

func (c *control) Error() Error {
	return newErrorManager(c)
}

func (c *control) Form(form *Form) {
	for _, field := range form.fields {
		field.value = c.ctx.FormValue(field.name, "")
	}
}

func (c *control) JSON(value any) {
	c.response = &response{
		responseType: responseJson,
		json:         value,
	}
}

func (c *control) Languages() []Language {
	result := make([]Language, len(c.app.languages))
	i := 0
	for _, l := range c.app.languages {
		result[i] = l
		i++
	}
	return result
}

func (c *control) Meta() Meta {
	return c.meta
}

func (c *control) Path() string {
	return c.ctx.Route().Path
}

func (c *control) Redirect(path string) {
	c.response.setType(responseRedirect).setRedirect(path)
}

func (c *control) Root() string {
	return root()
}

func (c *control) Template(path string) {
	c.response = &response{
		responseType: responseTemplate,
		template:     path,
	}
}

func (c *control) Text(value string) {
	c.response = &response{
		responseType: responseText,
		text:         value,
	}
}
