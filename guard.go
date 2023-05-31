package framework

type GuardHandler = func(control GuardControl) bool

type Guard interface {
	Handler(GuardHandler) Guard
	Redirect(redirect string) Guard
}

type guard struct {
	handler  GuardHandler
	redirect string
	name     string
}

const (
	guardDefault = "default"
)

func (g *guard) Handler(handler GuardHandler) Guard {
	g.handler = handler
	return g
}

func (g *guard) Redirect(redirect string) Guard {
	g.redirect = redirect
	return g
}
