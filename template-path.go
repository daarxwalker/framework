package framework

import (
	"strings"
)

type templatePath struct {
	path       string
	sourceType string
	namespace  string
	global     bool
}

const (
	templateSourceGlobal     = "global"
	templateSourceModule     = "module"
	templateSourceController = "controller"
	templateSourceComponent  = "component"
	templatePathSeparator    = ":"
)

func newTemplatePath(path, sourceType, namespace string) templatePath {
	tp := templatePath{
		path:       path,
		sourceType: sourceType,
		namespace:  namespace,
		global:     strings.HasPrefix(path, templateSourceGlobal+templatePathSeparator),
	}
	if strings.Contains(path, templatePathSeparator) {
		tp = tp.parse(path)
	}
	return tp
}

func (tp templatePath) parse(path string) templatePath {
	parts := strings.Split(path, templatePathSeparator)
	if len(parts) < 2 {
		return tp
	}
	tp.sourceType = parts[0]
	tp.path = parts[1]
	return tp
}

func (tp templatePath) buildPath() string {
	if tp.global {
		return templateSourceGlobal + templatePathSeparator + tp.path
	}
	path := tp.namespace + templatePathSeparator + tp.path
	if tp.sourceType == templateSourceController || tp.sourceType == templateSourceModule {
		path = tp.sourceType + templatePathSeparator + path
	}
	return path
}
