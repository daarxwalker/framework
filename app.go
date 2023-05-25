package framework

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/iancoleman/strcase"
)

type App struct {
	config          *Config
	controllers     map[string]*appController
	router          *Router
	templateManager *templateManager
	services        map[string]*appService
}

func New(config ...*Config) *App {
	cfg := &Config{Port: 6000}
	if len(config) > 0 {
		cfg = config[0]
	}
	return &App{
		config:          cfg,
		controllers:     make(map[string]*appController),
		templateManager: newTemplateManager(),
		services:        make(map[string]*appService),
	}
}

func (a *App) Router() *Router {
	a.router = newRouter(a)
	return a.router
}

func (a *App) Controller(controller any) *App {
	rp := &reflectProvider{
		reflectType:  reflect.TypeOf(controller),
		reflectValue: reflect.ValueOf(controller),
	}
	name := rp.reflectType.Elem().Name()
	name = strings.Replace(name, "Controller", "", -1)
	name = strcase.ToLowerCamel(name)
	a.controllers[name] = &appController{
		provider: rp,
		name:     name,
	}
	return a
}

func (a *App) Service(service any, config ...*ServiceConfig) *App {
	cfg := new(ServiceConfig)
	if len(config) > 0 {
		cfg = config[0]
	}
	rp := &reflectProvider{
		reflectType:  reflect.TypeOf(service),
		reflectValue: reflect.ValueOf(service),
	}
	if rp.reflectType.Kind() != reflect.Ptr {
		log.Fatalf("Service [%s] missing pointer.", rp.reflectType.Name())
	}
	a.services[rp.reflectValue.Elem().Type().String()] = &appService{
		provider: rp,
		config:   cfg,
	}
	return a
}

func (a *App) Fly() {
	greenBg := color.New(color.BgGreen)
	a.templateManager.parse()
	a.router.build()
	if _, err := greenBg.Println(fmt.Sprintf("Flying on host: [localhost:%d]", a.config.Port)); err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(a.router.fiber.Listen(fmt.Sprintf(":%d", a.config.Port)))
}
