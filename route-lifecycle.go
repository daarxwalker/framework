package framework

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type routeLifecycle struct {
	app                *App
	injector           *injector
	controller         *appController
	route              *Route
	control            *control
	routeControl       *routeControl
	namespace          *namespace
	ctx                *fiber.Ctx
	templateComponents map[string]reflect.Value
}

func newRouteLifecycle(app *App, controller *appController, route *Route, ctx *fiber.Ctx) *routeLifecycle {
	l := &routeLifecycle{
		app:                app,
		route:              route,
		ctx:                ctx,
		templateComponents: make(map[string]reflect.Value),
	}
	l.createContextController(controller)
	return l
}

func (l *routeLifecycle) run() {
	l.createNamespace()
	l.createControl()
	l.autoinject()
	l.processComponents()
	l.callRouteMethod()
	l.render()
}

func (l *routeLifecycle) autoinject() {
	l.injector = newInjector(l.app, l.control)
	l.injector.autoinject(l.controller.provider.reflectValue)
}

func (l *routeLifecycle) createNamespace() {
	l.namespace = newNamespace()
	l.namespace.set(l.controller.name)
}

func (l *routeLifecycle) createControl() {
	l.control = newControl(l.ctx)
	l.routeControl = newRouteControl(l.control, l.controller.name)
}

func (l *routeLifecycle) callRouteMethod() {
	l.route.method.Call([]reflect.Value{reflect.ValueOf(l.routeControl)})
}

func (l *routeLifecycle) createContextController(controller *appController) {
	c := reflect.New(controller.provider.reflectType.Elem())
	l.injector.copy(c, controller.provider.reflectValue.Elem())
	controller.provider.reflectValue = c
	l.controller = controller
}

func (l *routeLifecycle) processComponents() {
	cp := newComponentsProcessor(
		l.control,
		l.controller.provider.reflectValue,
		l.controller.provider.reflectType,
		l.namespace.clone(),
	)
	cp.route()
	cp.process()
	l.templateComponents = cp.result
}

func (l *routeLifecycle) render() {
	rm := newRenderManager(l.app, l.controller, l.routeControl.response.template, templateRoute, l.templateComponents)
	rm.render()
	if !rm.isOk() {
		l.routeControl.response.setStatus(fiber.StatusInternalServerError).setError(rm.error)
		return
	}
	l.routeControl.response.setStatus(fiber.StatusOK).setHtml(rm.html)
}
