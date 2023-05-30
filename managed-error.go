package framework

type managedError struct {
	control *control
	error   error
	message string
	status  int
}
