package framework

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mailgun/raymond/v2"
)

type TemplatesManager interface {
	GlobalPath(path string)
}

type templatesManager struct {
	app        *App
	globalPath string
	root       string
	partials   map[string]string
	layouts    map[string]*raymond.Template
	components map[string]*raymond.Template
	routes     map[string]*raymond.Template
}

const (
	templateControllerPathSuffix = "_controller"
	templateComponentPathSuffix  = "_component"
	templateModulePathSuffix     = "_module"
	templateFileTypeSuffix       = ".hbs"
)

var (
	templateComponent = "component"
	templateLayout    = "layout"
	templatePartial   = "partial"
	templateRoute     = "route"
	templateGlobal    = "global"
)

var (
	templateComponentSuffix = "." + templateComponent + templateFileTypeSuffix
	templateLayoutSuffix    = "." + templateLayout + templateFileTypeSuffix
	templatePartialSuffix   = "." + templatePartial + templateFileTypeSuffix
	templateRouteSuffix     = "." + templateRoute + templateFileTypeSuffix
)

func newTemplatesManager(app *App) *templatesManager {
	return &templatesManager{
		app:        app,
		partials:   make(map[string]string),
		layouts:    make(map[string]*raymond.Template),
		components: make(map[string]*raymond.Template),
		routes:     make(map[string]*raymond.Template),
		root:       root(),
	}
}

func (t *templatesManager) GlobalPath(path string) {
	t.globalPath = path
}

func (t *templatesManager) parse() {
	t.registerHelpers()
	t.parseTemplates(moduleRootDir)
	t.parseTemplates(t.globalPath)
}

func (t *templatesManager) parseTemplates(dir string) {
	check(filepath.Walk(t.root+dir, func(path string, info fs.FileInfo, err error) error {
		if !strings.HasSuffix(path, templateFileTypeSuffix) {
			return nil
		}
		switch {
		case strings.HasSuffix(path, templateComponentSuffix):
			t.parseTemplate(path, templateComponent)
		case strings.HasSuffix(path, templateLayoutSuffix):
			t.parseTemplate(path, templateLayout)
		case strings.HasSuffix(path, templatePartialSuffix):
			t.parseTemplate(path, templatePartial)
		case strings.HasSuffix(path, templateRouteSuffix):
			t.parseTemplate(path, templateRoute)
		}
		return nil
	}))
}

func (t *templatesManager) parseTemplate(path, templateType string) {
	content, err := os.ReadFile(path)
	check(err)
	if templateType == templatePartial {
		t.partials[t.getTemplateKey(path, templateType)] = string(content)
		return
	}
	tmpl, err := raymond.Parse(string(content))
	check(err)
	switch templateType {
	case templateComponent:
		t.components[t.getTemplateKey(path, templateType)] = tmpl
	case templateLayout:
		t.layouts[t.getTemplateKey(path, templateType)] = tmpl
	case templateRoute:
		t.routes[t.getTemplateKey(path, templateType)] = tmpl
	}
}

func (t *templatesManager) getTemplateKey(path, templateType string) string {
	var namespace, name string
	// Namespace
	isModuleChild := strings.Contains(path, templateModulePathSuffix)
	isControllerChild := strings.Contains(path, templateControllerPathSuffix)
	isComponentChild := strings.Contains(path, templateComponentPathSuffix)
	if isControllerChild {
		namespace = t.getControllerName(path)
	}
	if isComponentChild {
		namespace = t.getComponentName(path)
	}
	if isModuleChild && !isControllerChild && !isComponentChild {
		namespace = t.getModuleName(path)
	}
	// Name
	pathParts := strings.Split(path, "/")
	name = pathParts[len(pathParts)-1]
	switch templateType {
	case templateComponent:
		name = strings.TrimSuffix(name, templateComponentSuffix)
	case templateLayout:
		name = strings.TrimSuffix(name, templateLayoutSuffix)
	case templatePartial:
		name = strings.TrimSuffix(name, templatePartialSuffix)
	case templateRoute:
		name = strings.TrimSuffix(name, templateRouteSuffix)
	}
	if len(namespace) == 0 {
		return fmt.Sprintf("%s:%s", templateGlobal, name)
	}
	key := fmt.Sprintf("%s:%s", namespace, name)
	if templateType != templateRoute {
		return key
	}
	if isModuleChild && !isControllerChild {
		key = templateSourceModule + templatePathSeparator + key
	}
	if isControllerChild {
		key = templateSourceController + templatePathSeparator + key
	}
	return key
}

func (t *templatesManager) getControllerName(path string) (result string) {
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.HasSuffix(part, templateControllerPathSuffix) {
			return strings.TrimSuffix(part, templateControllerPathSuffix)
		}
	}
	return result
}

func (t *templatesManager) getComponentName(path string) (result string) {
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.HasSuffix(part, templateComponentPathSuffix) {
			return strings.TrimSuffix(part, templateComponentPathSuffix)
		}
	}
	return result
}

func (t *templatesManager) getModuleName(path string) (result string) {
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.HasSuffix(part, templateModulePathSuffix) {
			return strings.TrimSuffix(part, templateModulePathSuffix)
		}
	}
	return result
}

func (t *templatesManager) registerHelpers() {
	raymond.RegisterHelper("slot", func(options *raymond.Options) raymond.SafeString {
		return slotTag
	})
	raymond.RegisterHelper("form", func(action string, options *raymond.Options) raymond.SafeString {
		return raymond.SafeString(fmt.Sprintf(`<form action="%s" method="post">%s</form>`, action, options.Fn()))
	})
	raymond.RegisterHelper("link", func(key string, options *raymond.Options) raymond.SafeString {
		link, ok := t.app.linksManager.links[key]
		if !ok {
			return ""
		}
		return raymond.SafeString(link)
	})
}
