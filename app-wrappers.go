package framework

import "reflect"

type reflectProvider struct {
	reflectType  reflect.Type
	reflectValue reflect.Value
}

type appController struct {
	provider *reflectProvider
	name     string
	isModule bool
	module   string
}

type appModule struct {
	provider *reflectProvider
	name     string
}

type appService struct {
	provider *reflectProvider
	config   *ServiceConfig
	module   string
	isModule bool
}
