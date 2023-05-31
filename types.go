package framework

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type Ctx = fiber.Ctx
type MiddlewareHandler = fiber.Handler

var (
	componentType  = reflect.TypeOf(Component{})
	controllerType = reflect.TypeOf(Controller{})
	errorType      = reflect.TypeOf(managedError{})
	moduleType     = reflect.TypeOf(Module{})
	routerType     = reflect.TypeOf(Router{})
	serviceType    = reflect.TypeOf(Service{})
)
