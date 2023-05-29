package framework

type Component struct {
	Control  ComponentControl
	Handle   map[string]string
	TFS      templateFileSystem
	template templatePath
	name     string
	control  *control
}

func (c *Component) Template(path string) {
	c.template = newTemplatePath(path, templateSourceModule, c.control.module)
}
