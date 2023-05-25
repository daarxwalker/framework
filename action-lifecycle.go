package framework

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type actionLifecycle struct {
	name               string
	app                *App
	injector           *injector
	controller         *appController
	control            *control
	actionControl      *actionControl
	namespace          *namespace
	ctx                *fiber.Ctx
	templateComponents map[string]reflect.Value
	actionMethod       reflect.Value
}

const (
	queryAction = "action"
)

func newActionLifecycle(app *App, controller *appController, ctx *fiber.Ctx) *actionLifecycle {
	l := &actionLifecycle{
		name:               ctx.Query(queryAction, ""),
		app:                app,
		ctx:                ctx,
		templateComponents: make(map[string]reflect.Value),
	}
	l.createContextController(controller)
	return l
}

func (l *actionLifecycle) run() {
	if ok := l.verifyAction(); !ok {
		l.actionControl.response.setStatus(fiber.StatusBadRequest).setError(errors.New("unknown action"))
		return
	}
	l.createNamespace()
	l.createControl()
	l.autoinject()
	l.processComponents()
	l.callActionMethod()
}

func (l *actionLifecycle) verifyAction() bool {
	if len(l.name) == 0 {
		return false
	}
	return true
}

func (l *actionLifecycle) autoinject() {
	l.injector = newInjector(l.app, l.control)
	l.injector.autoinject(l.controller.provider.reflectValue)
}

func (l *actionLifecycle) createNamespace() {
	l.namespace = newNamespace()
	l.namespace.set(l.controller.name)
}

func (l *actionLifecycle) createControl() {
	l.control = newControl(l.ctx)
	l.actionControl = newActionControl(l.control)
}

func (l *actionLifecycle) callActionMethod() {
	if !l.actionMethod.IsValid() {
		return
	}
	l.actionMethod.Call([]reflect.Value{})
}

func (l *actionLifecycle) createContextController(controller *appController) {
	c := reflect.New(controller.provider.reflectType.Elem())
	l.injector.copy(c, controller.provider.reflectValue.Elem())
	controller.provider.reflectValue = c
	l.controller = controller
}

func (l *actionLifecycle) processComponents() {
	cp := newComponentsProcessor(
		l.control,
		l.controller.provider.reflectValue,
		l.controller.provider.reflectType,
		l.namespace.clone(),
	)
	cp.action(l.actionControl, l.createActionPath())
	cp.process()
	l.templateComponents = cp.result
	l.actionMethod = cp.actionMethod
}

func (l *actionLifecycle) createActionPath() []string {
	return strings.Split(l.name, "-")
}
