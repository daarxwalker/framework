package framework

type RouteControl interface {
	Text(value string) RouteControl
	Render(template string) RouteControl
	JSON(value any) RouteControl
}

type routeControl struct {
	control        *control
	response       *response
	controllerName string
}

func newRouteControl(control *control, controllerName string) *routeControl {
	return &routeControl{
		control:        control,
		controllerName: controllerName,
		response:       new(response),
	}
}

func (r *routeControl) Text(value string) RouteControl {
	r.response = &response{
		responseType: responseText,
		text:         value,
	}
	return r
}

func (r *routeControl) Render(template string) RouteControl {
	r.response = &response{
		responseType: responseTemplate,
		template:     template,
	}
	return r
}

func (r *routeControl) JSON(value any) RouteControl {
	r.response = &response{
		responseType: responseJson,
		json:         value,
	}
	return r
}
