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
	template    templatePath
	renderType  string
	error       error
	html        string
	layout      string
	contextData map[string]any
	components  map[string]reflect.Value
}

const (
	slotTag = "<slot/>"
)

func newRenderManager(app *App, controller *appController, template templatePath, renderType string, components map[string]reflect.Value) *renderManager {
	return &renderManager{
		app:         app,
		controller:  controller,
		template:    template,
		renderType:  renderType,
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
	tmpl, ok := m.getTemplates()[m.getTemplateKey()]
	if !ok {
		m.error = errors.New(fmt.Sprintf("%s template [%s/%s] not found.", m.renderType, m.controller.name, m.template.path))
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
		tmpl, ok := m.app.templatesManager.components[bc.template.buildPath()]
		if !ok {
			m.error = errors.New(fmt.Sprintf("component template [%s/%s] not found.", m.controller.name, m.template.path))
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
	if m.template.sourceType != templateSourceModule && m.template.sourceType != templateSourceController {
		return fmt.Sprintf("%s:%s", m.template.namespace, m.template.path)
	}
	return fmt.Sprintf("%s:%s:%s", m.template.sourceType, m.template.namespace, m.template.path)
}

func (m *renderManager) getTemplates() map[string]*raymond.Template {
	switch m.renderType {
	case templateComponent:
		return m.app.templatesManager.components
	case templateRoute:
		return m.app.templatesManager.routes
	default:
		return make(map[string]*raymond.Template)
	}
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
