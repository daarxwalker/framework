package framework

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

type componentLifecycle struct {
	baseComponent *Component
	component     reflect.Value
	componentType reflect.Type
	control       *control
	name          string
	namespace     *namespace
}

const (
	componentHandleMethod = "Handle"
	componentCreateMethod = "Create"
	componentRenderMethod = "Render"
)

func newComponentLifecycle(control *control, c reflect.Value, ct reflect.Type, name string, namespace *namespace) *componentLifecycle {
	return &componentLifecycle{
		control:       control,
		name:          name,
		namespace:     namespace,
		component:     c,
		componentType: ct,
		baseComponent: c.Elem().FieldByName(componentType.Name()).Interface().(*Component),
	}
}

func (l *componentLifecycle) createActions() {
	bc := l.component.Elem().FieldByName(componentType.Name()).Interface().(*Component)
	methodsLen := l.component.NumMethod()
	for i := 0; i < methodsLen; i++ {
		methodType := l.componentType.Method(i)
		if !strings.HasPrefix(methodType.Name, componentHandleMethod) {
			continue
		}
		bc.Handle[strcase.ToLowerCamel(strings.TrimPrefix(methodType.Name, componentHandleMethod))] = l.createActionUrl(l.namespace.create(methodType.Name))
	}
}

func (l *componentLifecycle) createActionUrl(actionName string) string {
	return fmt.Sprintf("%s?action=%s", l.control.Path(), actionName)
}

func (l *componentLifecycle) run() *componentLifecycle {
	l.createActions()
	l.runInitMethod()
	l.runCreateMethods()
	l.runRenderMethod()
	return l
}

func (l *componentLifecycle) runCreateMethods() {
	methodsLen := l.component.NumMethod()
	for i := 0; i < methodsLen; i++ {
		method := l.component.Method(i)
		methodType := l.componentType.Method(i)
		if !strings.HasPrefix(methodType.Name, componentCreateMethod) {
			continue
		}
		method.Call([]reflect.Value{})
	}
}

func (l *componentLifecycle) runInitMethod() {
	initMethod := l.component.MethodByName(initMethod)
	if !initMethod.IsValid() {
		return
	}
	initMethod.Call([]reflect.Value{})
}

func (l *componentLifecycle) runRenderMethod() {
	renderMethod := l.component.MethodByName(componentRenderMethod)
	if !renderMethod.IsValid() {
		return
	}
	renderMethod.Call([]reflect.Value{})
}
