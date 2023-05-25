package framework

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

type componentLifecycle struct {
	control       *control
	namespace     *namespace
	name          string
	baseComponent *Component
	component     reflect.Value
	componentType reflect.Type
}

const (
	componentInitMethod     = "Init"
	componentHandleMethod   = "Handle"
	componentCreateMethod   = "Create"
	componentTemplateMethod = "Template"
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

func (l *componentLifecycle) run() *componentLifecycle {
	l.createActions()
	l.runInitMethod()
	l.runCreateMethods()
	l.runTemplateMethod()
	return l
}

func (l *componentLifecycle) runInitMethod() {
	initMethod := l.component.MethodByName(componentInitMethod)
	if !initMethod.IsValid() {
		return
	}
	initMethod.Call([]reflect.Value{})
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

func (l *componentLifecycle) runTemplateMethod() {
	templateMethod := l.component.MethodByName(componentTemplateMethod)
	if !templateMethod.IsValid() {
		return
	}
	result := templateMethod.Call([]reflect.Value{})
	if len(result) == 0 {
		return
	}
	if len(l.name) == 0 {
		l.baseComponent.Template = result[0].String()
		return
	}
	l.baseComponent.Template = l.name + ":" + result[0].String()
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
