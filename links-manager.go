package framework

import "fmt"

type linksManager struct {
	app   *App
	links map[string]string
}

func newLinksManager(app *App) *linksManager {
	return &linksManager{
		app:   app,
		links: make(map[string]string),
	}
}

func (m *linksManager) build() {
	for _, route := range m.app.router.routes {
		m.add(route)
	}
}

func (m *linksManager) add(route *Route) {
	m.links[m.getLinkKey(route)] = m.getLinkPath(route)
}

func (m *linksManager) getLinkKey(route *Route) string {
	if route.isModule {
		return fmt.Sprintf("%s:%s:%s", route.module, route.controllerName, route.controllerMethod)
	}
	return fmt.Sprintf("%s:%s", route.controllerName, route.controllerMethod)
}

func (m *linksManager) getLinkPath(route *Route) string {
	return route.path
}
