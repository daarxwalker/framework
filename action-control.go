package framework

type actionControl struct {
	control  *control
	response *response
}

func newActionControl(control *control) *actionControl {
	return &actionControl{
		control:  control,
		response: new(response),
	}
}

func (c *actionControl) Form(dst any) {

}

func (c *actionControl) Render(name ...string) {
}

func (c *actionControl) Redirect(path string) {
	c.response.setType(responseRedirect).setRedirect(path)
}
