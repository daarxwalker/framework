package framework

import "reflect"

var (
	componentType  = reflect.TypeOf(Component{})
	controllerType = reflect.TypeOf(Controller{})
	moduleType     = reflect.TypeOf(Module{})
	routerType     = reflect.TypeOf(Router{})
	serviceType    = reflect.TypeOf(Service{})
)
