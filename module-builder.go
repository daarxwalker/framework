package framework

import (
	"reflect"
)

type moduleBuilder struct {
	app *App
}

const (
	moduleControllersMethod = "Controllers"
	moduleServicesMethod    = "Services"
	moduleRoutesMethod      = "Routes"
)

func newModuleBuilder(app *App) *moduleBuilder {
	return &moduleBuilder{app: app}
}

func (b *moduleBuilder) build() {
	for _, module := range b.app.modules {
		b.injectModule(module)
		b.runMethods(module)
		b.copyRoutesToMainRouter(module)
	}
}

func (b *moduleBuilder) injectModule(module *appModule) {
	router := newRouter(b.app)
	router.module = true
	router.moduleName = module.name
	m := b.getModule(module)
	if !m.IsValid() {
		return
	}
	m.Set(
		reflect.ValueOf(
			&Module{
				app:    b.app,
				Router: router,
				name:   module.name,
			},
		),
	)
}

func (b *moduleBuilder) runMethods(module *appModule) {
	b.runMethod(module, initMethod)
	b.runMethod(module, moduleControllersMethod)
	b.runMethod(module, moduleServicesMethod)
	b.runMethod(module, moduleRoutesMethod)
}

func (b *moduleBuilder) getModule(module *appModule) reflect.Value {
	m := module.provider.reflectValue
	if m.Kind() == reflect.Ptr {
		m = m.Elem()
	}
	result := m.FieldByName(moduleType.Name())
	if !result.IsValid() {
		return reflect.Value{}
	}
	return result
}

func (b *moduleBuilder) getRouter(module reflect.Value) *Router {
	if !module.IsValid() {
		return nil
	}
	if module.Kind() == reflect.Ptr {
		module = module.Elem()
	}
	router := module.FieldByName(routerType.Name())
	if !router.IsValid() {
		return nil
	}
	return router.Interface().(*Router)
}

func (b *moduleBuilder) runMethod(module *appModule, methodName string) {
	method := module.provider.reflectValue.MethodByName(methodName)
	if !method.IsValid() {
		return
	}
	method.Call([]reflect.Value{})
}

func (b *moduleBuilder) copyRoutesToMainRouter(module *appModule) {
	m := b.getModule(module)
	if !m.IsValid() {
		return
	}
	moduleRouter := b.getRouter(m)
	if moduleRouter == nil {
		return
	}
	b.app.router.routes = append(b.app.router.routes, moduleRouter.routes...)
}
