package framework

type appService struct {
	provider *reflectProvider
	config   *ServiceConfig
	module   string
	isModule bool
}
