package framework

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/iancoleman/strcase"
)

type Router struct {
	app        *App
	fiber      *fiber.App
	routes     []*Route
	module     bool
	moduleName string
	basePath   string
}

const (
	controllerRenderMethodPrefix = "Render"
	controllerDefaultRouteName   = "default"
	langParam                    = "lang"
)

func newRouter(app *App) *Router {
	return &Router{
		fiber: fiber.New(fiber.Config{
			DisableStartupMessage: true,
		}),
		app: app,
	}
}

func (r *Router) Group(path string) *Router {
	return &Router{
		fiber:    r.fiber,
		app:      r.app,
		basePath: path,
	}
}

func (r *Router) Add(path string) *Route {
	rr := &Route{path: r.basePath + path, module: r.moduleName, isModule: len(r.moduleName) > 0}
	r.routes = append(r.routes, rr)
	return rr
}

func (r *Router) build() {
	r.registerMiddlewares()
	r.buildControllers()
	r.createWildcardRoute()
}

func (r *Router) buildControllers() {
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
			if strings.HasPrefix(methodName, controllerRenderMethodPrefix) {
				r.buildRoute(controller, methodName)
				r.buildActions(controller, methodName)
			}
		}
	}
}

func (r *Router) buildRoute(controller *appController, renderMethodName string) {
	route := r.getRoute(renderMethodName)
	if route == nil {
		return
	}
	r.fiber.Get(r.modifyRoutePath(route.path), func(ctx *fiber.Ctx) error {
		l := newLifecycle(r.app, controller, ctx, renderMethodName, route)
		l.route()
		l.run()
		return r.buildResponse(ctx, l.control.response)
	})
}

func (r *Router) buildActions(controller *appController, renderMethodName string) {
	route := r.getRoute(renderMethodName)
	if route == nil {
		return
	}
	r.fiber.Post(r.modifyRoutePath(route.path), func(ctx *fiber.Ctx) error {
		l := newLifecycle(r.app, controller, ctx, renderMethodName, route)
		l.action()
		l.run()
		return r.buildResponse(ctx, l.control.response)
	})
}

func (r *Router) createWildcardRoute() {
	r.fiber.Get("*", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).SendString("not found")
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
	routeName := controllerDefaultRouteName
	if len(methodName) > len(controllerRenderMethodPrefix) {
		routeName = strcase.ToLowerCamel(strings.Replace(methodName, controllerRenderMethodPrefix, "", -1))
	}
	var resolvedRoute *Route
	for _, route := range r.routes {
		if route.controllerMethod != routeName {
			continue
		}
		resolvedRoute = route
	}
	return resolvedRoute
}

func (r *Router) modifyRoutePath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !r.app.i18n {
		return path
	}
	path = "/:" + langParam + path
	return path
}

func (r *Router) registerMiddlewares() {
	r.fiber.Use(middlewareError(r.app))
	r.fiber.Use(middlewareSlash(r.app))
	if r.app.i18n {
		r.fiber.Use(middlewareI18nMissingLang(r.app))
		r.fiber.Use("/:"+langParam+"/*", middlewareI18nVerifyLang(r.app))
	}
}
