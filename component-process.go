package framework

import (
	"reflect"
	"strings"
)

type componentProcess struct {
	control       *control
	namespace     *namespace
	component     reflect.Value
	typeComponent reflect.StructField
	processType   string
	actionPath    []string
	actionMethod  reflect.Value
}

const (
	componentTypeNameSuffix = "_component"
)

func (cp *componentProcess) createNamespace() {
	cp.namespace.set(cp.typeComponent.Name)
	cp.namespace = cp.namespace.clone()
}

func (cp *componentProcess) getFullname() string {
	var result string
	name := cp.component.Type()
	if name.Kind() == reflect.Ptr {
		name = name.Elem()
	}
	if strings.Contains(name.String(), componentTypeNameSuffix) {
		result = name.String()[:strings.Index(name.String(), componentTypeNameSuffix)]
	}
	return result
}

func (cp *componentProcess) getKind() reflect.Kind {
	return cp.component.Type().Kind()
}

func (cp *componentProcess) prepare() {
	cp.createNamespace()
}

func (cp *componentProcess) run() bool {
	if cp.processType == processAction && !cp.shouldActionProcessRun() {
		return false
	}
	switch cp.processType {
	case processAction:
		cp.runActionComponentProcess()
	case processRoute:
		cp.runRouteComponentProcess()
	}
	return true
}

func (cp *componentProcess) runActionComponentProcess() {
	if len(cp.namespace.names) < len(cp.actionPath)-1 {
		return
	}
	newComponentLifecycle(cp.control, cp.component, cp.typeComponent.Type, cp.getFullname(), cp.namespace.clone()).run()
	cp.setActionMethod()
}

func (cp *componentProcess) runRouteComponentProcess() {
	newComponentLifecycle(cp.control, cp.component, cp.typeComponent.Type, cp.getFullname(), cp.namespace.clone()).run()
}

func (cp *componentProcess) setActionMethod() {
	methodName := cp.actionPath[len(cp.actionPath)-1]
	cp.actionMethod = cp.component.MethodByName(methodName)
}

func (cp *componentProcess) shouldActionProcessRun() bool {
	if cp.processType == processRoute {
		return true
	}
	return cp.verifyActionNamespace()
}

func (cp *componentProcess) shouldProcess() bool {
	kind := cp.getKind()
	if kind != reflect.Ptr {
		return false
	}
	if cp.component.IsNil() {
		return false
	}
	if !cp.component.Elem().FieldByName(componentType.Name()).IsValid() {
		return false
	}
	if cp.component.Elem().NumField() == 0 {
		return false
	}
	return true
}

func (cp *componentProcess) verifyActionNamespace() bool {
	if len(cp.namespace.names) > len(cp.actionPath) {
		return false
	}
	for i, item := range cp.namespace.names {
		if item != cp.actionPath[i] {
			return false
		}
	}
	return true
}
