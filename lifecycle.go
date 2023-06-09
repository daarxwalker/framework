package framework

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"

	"dd"
)

type lifecycle struct {
	app                *App
	injector           *injector
	controller         *appController
	control            *control
	namespace          *namespace
	ctx                *fiber.Ctx
	templateComponents map[string]reflect.Value
	renderMethodName   string
	guard              string
	module             string
	lifecycleType      string
	actionName         string
	actionMethod       reflect.Value
}

const (
	lifecycleRoute  = "route"
	lifecycleAction = "action"
	queryAction     = "action"
)

func newLifecycle(app *App, controller *appController, ctx *fiber.Ctx, renderMethodName string, route *Route) *lifecycle {
	l := &lifecycle{
		app:                app,
		renderMethodName:   renderMethodName,
		module:             route.module,
		guard:              route.guard,
		ctx:                ctx,
		actionName:         ctx.Query(queryAction, ""),
		templateComponents: make(map[string]reflect.Value),
	}
	l.createContextController(controller)
	return l
}

func (l *lifecycle) route() {
	l.lifecycleType = lifecycleRoute
}

func (l *lifecycle) action() {
	l.lifecycleType = lifecycleAction
}

func (l *lifecycle) run() {
	l.beforeInject()
	dd.Print(l.validateGuard())
	if ok := l.validateGuard(); !ok {
		return
	}
	l.autoinject()
	l.beforeRender()
	l.render()
}

func (l *lifecycle) autoinject() {
	l.injector = newInjector(l.app, l.control)
	l.injector.autoinject(l.controller.provider.reflectValue, l.controller.provider.reflectType)
}

func (l *lifecycle) createNamespace() {
	l.namespace = newNamespace()
	l.namespace.set(l.controller.name)
}

func (l *lifecycle) createControl() {
	l.control = newControl(l.app, l.ctx, l.controller, l.module)
}

func (l *lifecycle) callInitMethod() {
	method := l.controller.provider.reflectValue.MethodByName(initMethod)
	if !method.IsValid() {
		return
	}
	method.Call([]reflect.Value{})
}

func (l *lifecycle) callRenderMethod() {
	method := l.controller.provider.reflectValue.MethodByName(l.renderMethodName)
	if !method.IsValid() {
		return
	}
	method.Call([]reflect.Value{})
}

func (l *lifecycle) createContextController(controller *appController) {
	c := reflect.New(controller.provider.reflectType.Elem())
	l.injector.copy(c, controller.provider.reflectValue.Elem())
	controller.provider.reflectValue = c
	l.controller = controller
}

func (l *lifecycle) processComponents() {
	cp := newComponentsProcessor(
		l.control,
		l.controller.provider.reflectValue,
		l.controller.provider.reflectType,
		l.namespace.clone(),
	)
	switch l.lifecycleType {
	case lifecycleRoute:
		cp.route()
	case lifecycleAction:
		cp.action(l.createActionPath())
	}
	cp.process()
	l.templateComponents = cp.result
	if l.lifecycleType == lifecycleAction {
		l.actionMethod = cp.actionMethod
	}
}

func (l *lifecycle) callActionMethod() {
	if l.lifecycleType != lifecycleAction {
		return
	}
	l.actionMethod.Call([]reflect.Value{})
}

func (l *lifecycle) beforeInject() {
	l.createNamespace()
	l.createControl()
}

func (l *lifecycle) beforeRender() {
	l.callInitMethod()
	l.processComponents()
	l.callActionMethod()
}

func (l *lifecycle) render() {
	if len(l.control.response.responseType) != 0 {
		return
	}
	l.callRenderMethod()
	rm := newRenderManager(l.app, l.control, l.controller, l.control.response.template, l.templateComponents)
	rm.render()
	if !rm.isOk() {
		l.control.response.setStatus(fiber.StatusInternalServerError).setError(rm.error)
		return
	}
	l.control.response.setStatus(fiber.StatusOK).setHtml(rm.html)
}

func (l *lifecycle) verifyAction() bool {
	if len(l.actionName) == 0 {
		return false
	}
	return true
}

func (l *lifecycle) createActionPath() []string {
	return strings.Split(l.actionName, "-")
}

func (l *lifecycle) validateGuard() bool {
	if len(l.guard) == 0 {
		return true
	}
	guard, ok := l.app.guards[l.guard]
	if !ok {
		l.control.response.setStatus(fiber.StatusForbidden).setError(errors.New("unknown guard - prevently forbidden"))
		return false
	}
	return guard.handler(l.control)
}
