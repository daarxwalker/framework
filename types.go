package framework

import "reflect"

var (
	componentType  = reflect.TypeOf(Component{})
	controllerType = reflect.TypeOf(Controller{})
	serviceType    = reflect.TypeOf(Service{})
)
