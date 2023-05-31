package framework

type Component struct {
	Control ComponentControl
	Error   func() Error
	Handle  map[string]string

	control  *control
	name     string
	template string
}

func (c *Component) Template(path string) {
	c.template = path
}
