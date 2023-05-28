package framework

import (
	"reflect"

	"github.com/iancoleman/strcase"
)

type componentsProcessor struct {
	control        *control
	namespace      *namespace
	targetValue    reflect.Value
	targetType     reflect.Type
	targetValuePtr reflect.Value
	targetTypePtr  reflect.Type
	processType    string
	actionPath     []string
	actionMethod   reflect.Value
	result         map[string]reflect.Value
}

const (
	processRoute  = "route"
	processAction = "action"
)

func newComponentsProcessor(control *control, targetValue reflect.Value, targetType reflect.Type, namespace *namespace) *componentsProcessor {
	cp := &componentsProcessor{
		control:        control,
		namespace:      namespace,
		targetValue:    targetValue.Elem(),
		targetType:     targetType.Elem(),
		targetValuePtr: targetValue,
		targetTypePtr:  targetType,
		result:         make(map[string]reflect.Value),
	}
	return cp
}

func (cp *componentsProcessor) action(path []string) {
	cp.processType = processAction
	cp.actionPath = path
}

func (cp *componentsProcessor) route() {
	cp.processType = processRoute
}

func (cp *componentsProcessor) process() {
	cp.processComponents(cp.targetValuePtr, cp.targetTypePtr, cp.namespace.clone())
}

func (cp *componentsProcessor) processComponents(targetValue reflect.Value,
	targetType reflect.Type, namespace *namespace) {
	fieldsLen := targetValue.Elem().NumField()
	if fieldsLen == 0 {
		return
	}
	for index := 0; index < fieldsLen; index++ {
		p := &componentProcess{
			control:       cp.control,
			processType:   cp.processType,
			actionPath:    cp.actionPath,
			namespace:     namespace.clone(),
			component:     targetValue.Elem().Field(index),
			typeComponent: targetType.Elem().Field(index),
		}
		if !p.shouldProcess() {
			continue
		}
		p.prepare()
		p.run()
		cp.updateResult(p.component)
		if cp.processType == processAction && p.actionMethod.IsValid() {
			cp.actionMethod = p.actionMethod
		}
		cp.processComponents(p.component, p.typeComponent.Type, p.namespace.clone())
	}
}

func (cp *componentsProcessor) updateResult(target reflect.Value) {
	if !target.Elem().FieldByName(componentType.Name()).IsValid() {
		return
	}
	cp.result[strcase.ToKebab(target.Elem().Type().Name())] = target
}
