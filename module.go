package framework

type Module struct {
	Router *Router
	app    *App
	name   string
}

func (m *Module) Controller(controller any) *Module {
	m.app.controller(controller, m.name)
	return m
}

func (m *Module) Service(service any, config ...*ServiceConfig) *Module {
	cfg := new(ServiceConfig)
	if len(config) > 0 {
		cfg = config[0]
	}
	cfg.module = m.name
	cfg.isModule = len(m.name) > 0
	m.app.Service(service, cfg)
	return m
}
