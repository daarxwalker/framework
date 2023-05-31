package framework

import "reflect"

type Route struct {
	path             string
	controllerName   string
	controllerMethod string
	method           reflect.Value
	module           string
	isModule         bool
	guard            string
}

func (r *Route) Controller(name string) *Route {
	r.controllerName = name
	return r
}

func (r *Route) Method(name string) *Route {
	r.controllerMethod = name
	return r
}

func (r *Route) Guard(name ...string) *Route {
	guardName := guardDefault
	if len(name) > 0 {
		guardName = name[0]
	}
	r.guard = guardName
	return r
}
