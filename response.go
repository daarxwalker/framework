package framework

type response struct {
	responseType string
	bytes        []byte
	json         any
	text         string
	template     string
	html         string
	error        error
	status       int
	redirect     string
}

const (
	responseText     = "text"
	responseTemplate = "template"
	responseJson     = "json"
	responseRedirect = "redirect"
)

func (r *response) setType(responseType string) *response {
	r.responseType = responseType
	return r
}

func (r *response) setError(err error) *response {
	r.error = err
	return r
}

func (r *response) setStatus(status int) *response {
	r.status = status
	return r
}

func (r *response) setHtml(html string) *response {
	r.html = html
	return r
}

func (r *response) setRedirect(redirect string) *response {
	r.redirect = redirect
	return r
}
