package framework

type formFieldRule struct {
	name string
	min  int
	max  int
}

const (
	formFieldRuleRequired = "required"
)
