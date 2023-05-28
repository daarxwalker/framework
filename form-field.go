package framework

type FormFieldBuilder interface {
	Type(fieldType string) *FormField
	Value(value any) *FormField
	Required() *FormField
}

type FormFieldGetter interface {
	String() string
	Bool() bool
}

type FormField struct {
	fieldType string
	name      string
	element   string
	rules     []*formFieldRule
	value     any
}

func (f *FormField) Type(fieldType string) *FormField {
	f.fieldType = fieldType
	return f
}

func (f *FormField) Value(value any) *FormField {
	f.value = value
	return f
}

func (f *FormField) Required() *FormField {
	rule := &formFieldRule{name: formFieldRuleRequired}
	f.rules = append(f.rules, rule)
	return f
}

func (f *FormField) String() string {
	return f.value.(string)
}

func (f *FormField) Bool() bool {
	return f.value.(string) == "true"
}
