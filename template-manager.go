package framework

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mailgun/raymond/v2"
)

type templateManager struct {
	app        *App
	partials   map[string]string
	layouts    map[string]*raymond.Template
	components map[string]*raymond.Template
	routes     map[string]*raymond.Template
}

const (
	templateControllerPathSuffix = "_controller"
	templateComponentPathSuffix  = "_component"
	templateDomainPathSuffix     = "_domain"
	templateFileTypeSuffix       = ".hbs"
)

var (
	templateComponent = "component"
	templateLayout    = "layout"
	templatePartial   = "partial"
	templateRoute     = "route"
)

var (
	templateComponentSuffix = "." + templateComponent + templateFileTypeSuffix
	templateLayoutSuffix    = "." + templateLayout + templateFileTypeSuffix
	templatePartialSuffix   = "." + templatePartial + templateFileTypeSuffix
	templateRouteSuffix     = "." + templateRoute + templateFileTypeSuffix
)

func newTemplateManager(app *App) *templateManager {
	return &templateManager{
		app:        app,
		partials:   make(map[string]string),
		layouts:    make(map[string]*raymond.Template),
		components: make(map[string]*raymond.Template),
		routes:     make(map[string]*raymond.Template),
	}
}

func (t *templateManager) parse() {
	t.registerHelpers()
	t.parseTemplates(domainRootDir)
	t.parseTemplates(templateRootDir)
}

func (t *templateManager) parseTemplates(dir string) {
	check(filepath.Walk(root()+dir, func(path string, info fs.FileInfo, err error) error {
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

func (t *templateManager) parseTemplate(path, templateType string) {
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

func (t *templateManager) getTemplateKey(path, templateType string) string {
	var namespace, name string
	// Namespace
	isDomainChild := strings.Contains(path, templateDomainPathSuffix)
	isControllerChild := strings.Contains(path, templateControllerPathSuffix)
	isComponentChild := strings.Contains(path, templateComponentPathSuffix)
	if isControllerChild {
		namespace = t.getControllerName(path)
	}
	if isComponentChild {
		namespace = t.getComponentName(path)
	}
	if isDomainChild && !isControllerChild && !isComponentChild {
		namespace = t.getDomainName(path)
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
		return name
	}
	return fmt.Sprintf("%s:%s", namespace, name)
}

func (t *templateManager) getControllerName(path string) (result string) {
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.HasSuffix(part, templateControllerPathSuffix) {
			return strings.TrimSuffix(part, templateControllerPathSuffix)
		}
	}
	return result
}

func (t *templateManager) getComponentName(path string) (result string) {
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.HasSuffix(part, templateComponentPathSuffix) {
			return strings.TrimSuffix(part, templateComponentPathSuffix)
		}
	}
	return result
}

func (t *templateManager) getDomainName(path string) (result string) {
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.HasSuffix(part, templateDomainPathSuffix) {
			return strings.TrimSuffix(part, templateDomainPathSuffix)
		}
	}
	return result
}

func (t *templateManager) registerHelpers() {
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
