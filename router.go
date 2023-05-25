package framework

import (
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/iancoleman/strcase"

	"dd"
)

type Router struct {
	app    *App
	fiber  *fiber.App
	routes []*Route
}

func newRouter(app *App) *Router {
	return &Router{
		fiber: fiber.New(fiber.Config{
			DisableStartupMessage: true,
		}),
		app: app,
	}
}

func (r *Router) Add(path string) *Route {
	rr := &Route{path: path}
	r.routes = append(r.routes, rr)
	return rr
}

func (r *Router) build() {
	for controllerName, controller := range r.app.controllers {
		ct := controller.provider.reflectType
		cv := controller.provider.reflectValue
		routes := make([]*Route, 0)
		for _, route := range r.routes {
			if route.controllerName != controllerName {
				continue
			}
			routes = append(routes, route)
		}
		methodsLen := cv.NumMethod()
		if methodsLen == 0 {
			continue
		}
		for index := 0; index < methodsLen; index++ {
			methodName := ct.Method(index).Name
			method := cv.Method(index)
			if strings.HasSuffix(methodName, "Route") {
				r.buildRoute(controller, methodName, method)
				r.buildActions(controller, methodName)
			}
		}
	}
}

func (r *Router) buildRoute(controller *appController, methodName string, method reflect.Value) {
	route := r.getRoute(methodName)
	if route == nil {
		return
	}
	route.method = method
	r.fiber.Get(route.path, func(ctx *fiber.Ctx) error {
		dd.Print("GET: ", ctx)
		l := newRouteLifecycle(r.app, controller, route, ctx)
		l.run()
		return r.buildResponse(ctx, l.routeControl.response)
	})
}

func (r *Router) buildActions(controller *appController, methodName string) {
	route := r.getRoute(methodName)
	if route == nil {
		return
	}
	r.fiber.Post(route.path, func(ctx *fiber.Ctx) error {
		dd.Print("POST: ", ctx)
		l := newActionLifecycle(r.app, controller, ctx)
		l.run()
		return r.buildResponse(ctx, l.actionControl.response)
	})
}

func (r *Router) buildResponse(ctx *fiber.Ctx, response *response) error {
	switch response.responseType {
	case responseRedirect:
		return ctx.Redirect(response.redirect, fiber.StatusSeeOther)
	case responseTemplate:
		ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		if response.error != nil {
			return ctx.Status(response.status).SendString(response.error.Error())
		}
		return ctx.Status(response.status).SendString(response.html)
	case responseJson:
		return ctx.JSON(response.json)
	case responseText:
		return ctx.SendString(response.text)
	default:
		return ctx.Send(response.bytes)
	}
}

func (r *Router) getRoute(methodName string) *Route {
	routeName := strcase.ToLowerCamel(strings.Replace(methodName, "Route", "", -1))
	var resolvedRoute *Route
	for _, route := range r.routes {
		if route.controllerMethod != routeName {
			continue
		}
		resolvedRoute = route
	}
	return resolvedRoute
}
