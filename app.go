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
	config           *Config
	controllers      map[string]*appController
	i18n             bool
	languages        map[string]*language
	linksManager     *linksManager
	modules          map[string]*appModule
	moduleBuilder    *moduleBuilder
	router           *Router
	services         map[string]*appService
	templatesManager *templatesManager
}

const (
	controllerSuffix = "Controller"
	moduleSuffix     = "Module"
)

func New(config ...*Config) *App {
	cfg := getAppDefaultConfig()
	if len(config) > 0 {
		cfg = config[0]
	}
	app := &App{
		config:      cfg,
		controllers: make(map[string]*appController),
		languages:   make(map[string]*language),
		modules:     make(map[string]*appModule),
		services:    make(map[string]*appService),
	}
	app.templatesManager = newTemplatesManager(app)
	app.router = newRouter(app)
	app.moduleBuilder = newModuleBuilder(app)
	app.linksManager = newLinksManager(app)
	return app
}

func (a *App) Controller(controller any) {
	a.controller(controller, "")
}

func (a *App) Fly() {
	a.beforeFly()
	greenBg := color.New(color.BgGreen)
	if _, err := greenBg.Println(fmt.Sprintf("Flying on host: [localhost:%d]", a.config.Port)); err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(a.router.fiber.Listen(fmt.Sprintf(":%d", a.config.Port)))
}

func (a *App) I18N() {
	a.i18n = true
}

func (a *App) Language(code string, main ...bool) Language {
	lang := newLanguage(code, main...)
	a.languages[code] = lang
	return lang
}

func (a *App) Module(module any) *App {
	rp := &reflectProvider{
		reflectType:  reflect.TypeOf(module),
		reflectValue: reflect.ValueOf(module),
	}
	name := rp.reflectType.Elem().Name()
	name = strings.Replace(name, moduleSuffix, "", -1)
	name = strcase.ToLowerCamel(name)
	a.modules[name] = &appModule{
		provider: rp,
		name:     name,
	}
	return a
}

func (a *App) Router() *Router {
	return a.router
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
		module:   cfg.module,
		isModule: len(cfg.module) > 0,
	}
	return a
}

func (a *App) Templates() TemplatesManager {
	return a.templatesManager
}

func (a *App) beforeFly() {
	a.templatesManager.parse()
	a.moduleBuilder.build()
	a.router.build()
	a.linksManager.build()
}

func (a *App) controller(controller any, module string) *App {
	rp := &reflectProvider{
		reflectType:  reflect.TypeOf(controller),
		reflectValue: reflect.ValueOf(controller),
	}
	name := rp.reflectType.Elem().Name()
	name = strings.Replace(name, controllerSuffix, "", -1)
	name = strcase.ToLowerCamel(name)
	a.controllers[name] = &appController{
		provider: rp,
		name:     name,
		module:   module,
		isModule: len(module) > 0,
	}
	return a
}

func getAppDefaultConfig() *Config {
	return &Config{Port: 6000}
}
