package framework

type Form[T any] struct {
	fields []*FormField
	state  T
}

func CreateForm[T any](initialState ...T) *Form[T] {
	f := &Form[T]{}
	if len(initialState) > 0 {
		f.state = initialState[0]
	}
	return f
}

func (f *Form[T]) Field(name string) *FormField {
	field := &FormField{name: name}
	f.fields = append(f.fields, field)
	return field
}

func (f *Form[T]) Update(fn func(value T) T) *Form[T] {
	f.state = fn(f.state)
	return f
}

func (f *Form[T]) Set(value T) *Form[T] {
	f.state = value
	return f
}

func (f *Form[T]) Validate(value T) bool {
	return true
}
