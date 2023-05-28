package framework

type Form struct {
	fields []*FormField
}

func CreateForm() *Form {
	f := &Form{}
	return f
}

func (f *Form) Field(name string) FormFieldBuilder {
	field := &FormField{name: name}
	f.fields = append(f.fields, field)
	return field
}

func (f *Form) Validate() bool {
	return true
}

func (f *Form) Get(name string) FormFieldGetter {
	for _, field := range f.fields {
		if field.name == name {
			return field
		}
	}
	return &FormField{}
}
