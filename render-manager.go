package framework

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/mailgun/raymond/v2"
)

type renderManager struct {
	app         *App
	controller  *appController
	template    string
	error       error
	html        string
	layout      string
	contextData map[string]any
	components  map[string]reflect.Value
}

const (
	slotTag = "<slot/>"
)

func newRenderManager(app *App, controller *appController, template string, components map[string]reflect.Value) *renderManager {
	return &renderManager{
		app:         app,
		controller:  controller,
		template:    template,
		components:  components,
		contextData: make(map[string]any),
	}
}

func (m *renderManager) render() {
	m.includeControllerContext()
	m.renderComponents()
	m.renderTemplate()
}

func (m *renderManager) includeControllerContext() {
	m.contextData["controller"] = m.controller.provider.reflectValue
}

func (m *renderManager) renderTemplate() {
	layouts := m.getLayouts()
	tmpl, ok := m.app.templatesManager.templates[m.getTemplateKey()]
	if !ok {
		m.error = errors.New(fmt.Sprintf("template [%s] not found.", m.template))
		return
	}
	if len(m.layout) == 0 {
		m.layout = "default"
	}
	m.html, m.error = tmpl.Exec(m.contextData)
	layoutTmpl, ok := layouts[m.layout]
	if ok {
		m.contextData["router"] = func(options *raymond.Options) raymond.SafeString {
			routeHtml := m.html
			return raymond.SafeString(routeHtml)
		}
		m.html, m.error = layoutTmpl.Exec(m.contextData)
		return
	}
}

func (m *renderManager) renderComponents() {
	for name, component := range m.components {
		bc := component.Elem().FieldByName(componentType.Name()).Interface().(*Component)
		tmpl, ok := m.app.templatesManager.templates[bc.template]
		if !ok {
			m.error = errors.New(fmt.Sprintf("component template [%s] not found.", m.template))
			return
		}
		m.contextData[name] = func(options *raymond.Options) raymond.SafeString {
			m.overrideComponentParams(options.Hash(), component)
			html, err := tmpl.Exec(component.Interface())
			if err != nil {
				return raymond.SafeString(err.Error())
			}
			html = strings.ReplaceAll(html, slotTag, options.Fn())
			return raymond.SafeString(html)
		}
	}
}

func (m *renderManager) overrideComponentParams(hash map[string]any, component reflect.Value) {
	if component.Type().Kind() == reflect.Ptr {
		component = component.Elem()
	}
	fieldsLen := component.NumField()
	if fieldsLen == 0 {
		return
	}
	for name, value := range hash {
		field := component.FieldByName(strcase.ToCamel(name))
		if !field.IsValid() {
			continue
		}
		field.Set(reflect.ValueOf(value))
	}
}

func (m *renderManager) getTemplateKey() string {
	return m.template
}

func (m *renderManager) getPartials() map[string]string {
	return m.app.templatesManager.partials
}

func (m *renderManager) getLayouts() map[string]*raymond.Template {
	return m.app.templatesManager.layouts
}

func (m *renderManager) isOk() bool {
	return m.error == nil
}
