package framework

type FormField struct {
	fieldType string
	name      string
	element   string
	rules     []*formFieldRule
}

func (f *FormField) Type(fieldType string) *FormField {
	f.fieldType = fieldType
	return f
}

func (f *FormField) Required() *FormField {
	rule := &formFieldRule{name: formFieldRuleRequired}
	f.rules = append(f.rules, rule)
	return f
}
