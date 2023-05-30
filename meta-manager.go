package framework

type Meta interface {
	Description(description string) Meta
	Keywords(keywords string) Meta
	Title(title string) Meta
}

type metaManager struct {
	title       string
	description string
	keywords    string
}

func (m *metaManager) Title(title string) Meta {
	m.title = title
	return m
}

func (m *metaManager) Description(description string) Meta {
	m.description = description
	return m
}

func (m *metaManager) Keywords(keywords string) Meta {
	m.keywords = keywords
	return m
}

func (m *metaManager) getMap() map[string]any {
	return map[string]any{
		"title":       m.title,
		"description": m.description,
		"keywords":    m.keywords,
	}
}
