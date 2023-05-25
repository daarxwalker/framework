package framework

import "strings"

type namespace struct {
	names []string
}

func newNamespace() *namespace {
	return &namespace{}
}

func (n *namespace) set(name string) {
	n.names = append(n.names, name)
}

func (n *namespace) clone() *namespace {
	return &namespace{names: n.names}
}

func (n *namespace) get() string {
	return strings.Join(n.names, "-")
}

func (n *namespace) create(name string) string {
	return strings.Join(append(n.names, name), "-")
}
