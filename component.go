package framework

type Component struct {
	Control  ComponentControl
	Handle   map[string]string
	Template string
}
