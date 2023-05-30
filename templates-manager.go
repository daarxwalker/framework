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
	Path(path string)
}

type templatesManager struct {
	app       *App
	path      string
	root      string
	partials  map[string]string
	layouts   map[string]*raymond.Template
	templates map[string]*raymond.Template
}

const (
	templateFileTypeSuffix = ".hbs"
)

var (
	templateLayout  = "layout"
	templatePartial = "partial"
)

var (
	templatePartialSuffix = "." + templatePartial + templateFileTypeSuffix
	templateLayoutSuffix  = "." + templateLayout + templateFileTypeSuffix
)

func newTemplatesManager(app *App) *templatesManager {
	return &templatesManager{
		app:       app,
		partials:  make(map[string]string),
		layouts:   make(map[string]*raymond.Template),
		templates: make(map[string]*raymond.Template),
		root:      root(),
	}
}

func (t *templatesManager) Path(path string) {
	t.path = path
}

func (t *templatesManager) parse() {
	t.registerHelpers()
	t.parseTemplates()
}

func (t *templatesManager) parseTemplates() {
	check(filepath.Walk(t.root+t.path, func(path string, info fs.FileInfo, err error) error {
		if !strings.HasSuffix(path, templateFileTypeSuffix) {
			return nil
		}
		switch {
		case strings.HasSuffix(path, templatePartialSuffix):
			t.parseTemplate(path, templatePartial)
		case strings.HasSuffix(path, templateLayoutSuffix):
			t.parseTemplate(path, templateLayout)
		default:
			t.parseTemplate(path, "")
		}
		return nil
	}))
}

func (t *templatesManager) parseTemplate(path string, templateType string) {
	key := t.getTemplateKey(path, templateType)
	content, err := os.ReadFile(path)
	check(err)
	if templateType == templatePartial {
		t.partials[key] = string(content)
		return
	}
	tmpl, err := raymond.Parse(string(content))
	check(err)
	if templateType == templateLayout {
		t.layouts[key] = tmpl
		return
	}
	t.templates[key] = tmpl
}

func (t *templatesManager) getTemplateKey(path string, templateType string) string {
	var suffix string
	switch templateType {
	case templatePartial:
		suffix = templatePartialSuffix
	case templateLayout:
		suffix = templateLayoutSuffix
	default:
		suffix = templateFileTypeSuffix
	}
	path = strings.TrimPrefix(path, t.root+t.path)
	path = strings.TrimSuffix(path, suffix)
	path = strings.TrimPrefix(path, "/")
	return path
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
