package framework

import (
	"github.com/gofiber/fiber/v2"
)

type BaseControl interface {
	Root() string
	CurrentLang() string
	Languages() []Language
}

type ComponentControl interface {
	BaseControl
	Form(form *Form)
	Redirect(path string)
}

type RouteControl interface {
	BaseControl
	Text(value string)
	Template(template string)
	JSON(value any)
	Redirect(path string)
}

type ControllerControl interface {
	BaseControl
	Text(value string)
	Template(template string)
	JSON(value any)
	Redirect(path string)
}

type control struct {
	app        *App
	ctx        *fiber.Ctx
	response   *response
	controller *appController
	module     string
}

func newControl(app *App, ctx *fiber.Ctx, controller *appController, module string) *control {
	return &control{
		app:        app,
		ctx:        ctx,
		response:   new(response),
		controller: controller,
		module:     module,
	}
}

func (c *control) CurrentLang() string {
	return c.ctx.Params(langParam)
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

func (c *control) Root() string {
	return root()
}

func (c *control) Path() string {
	return c.ctx.Route().Path
}

func (c *control) Text(value string) {
	c.response = &response{
		responseType: responseText,
		text:         value,
	}
}

func (c *control) Template(path string) {
	c.response = &response{
		responseType: responseTemplate,
		template:     path,
	}
}

func (c *control) JSON(value any) {
	c.response = &response{
		responseType: responseJson,
		json:         value,
	}
}

func (c *control) Form(form *Form) {
	for _, field := range form.fields {
		field.value = c.ctx.FormValue(field.name, "")
	}
}

func (c *control) Redirect(path string) {
	c.response.setType(responseRedirect).setRedirect(path)
}
