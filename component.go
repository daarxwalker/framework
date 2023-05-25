package framework

type Component struct {
	control  *control
	Action   *actionControl
	Handle   map[string]string
	Template string
}

const (
	componentActionFieldKey = "Action"
)
