package framework

type Language interface {
	Code() string
	Locale(key, value string, params ...map[string]any) Language
	Main() bool
}

type language struct {
	code    string
	locales map[string]*locale
	main    bool
}

func newLanguage(code string, main ...bool) *language {
	l := &language{code: code, locales: make(map[string]*locale)}
	if len(main) > 0 {
		l.main = main[0]
	}
	return l
}

func (l *language) Code() string {
	return l.code
}

func (l *language) Locale(key, value string, params ...map[string]any) Language {
	localeParams := make(map[string]any)
	if len(params) > 0 {
		localeParams = params[0]
	}
	l.locales[key] = &locale{
		key:    key,
		value:  value,
		params: localeParams,
	}
	return l
}

func (l *language) Main() bool {
	return l.main
}
