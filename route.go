package framework

import "reflect"

type Route struct {
	path             string
	controllerName   string
	controllerMethod string
	method           reflect.Value
}

func (r *Route) Controller(name string) *Route {
	r.controllerName = name
	return r
}

func (r *Route) Method(name string) *Route {
	r.controllerMethod = name
	return r
}
