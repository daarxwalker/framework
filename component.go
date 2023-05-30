package framework

type Component struct {
	Control  ComponentControl
	Handle   map[string]string
	template string
	name     string
	control  *control
}

func (c *Component) Template(path string) {
	c.template = path
}
