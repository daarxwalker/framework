package framework

type Controller struct {
	Control ControllerControl
	Error   func() Error
}
