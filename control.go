package framework

import (
	"github.com/gofiber/fiber/v2"
)

type ComponentControl interface {
	Form(form *Form)
	Redirect(path string)
}

type RouteControl interface {
	Text(value string)
	Template(template string)
	JSON(value any)
	Redirect(path string)
}

type ControllerControl interface {
	Text(value string)
	Template(template string)
	JSON(value any)
	Redirect(path string)
}

type control struct {
	ctx        *fiber.Ctx
	response   *response
	controller *appController
}

func newControl(ctx *fiber.Ctx) *control {
	return &control{
		ctx:        ctx,
		response:   new(response),
		controller: new(appController),
	}
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

func (c *control) Template(template string) {
	c.response = &response{
		responseType: responseTemplate,
		template:     template,
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
