package framework

import (
	"log"
	"reflect"
)

type injector struct {
	app     *App
	control *control
}

func newInjector(app *App, control *control) *injector {
	return &injector{
		app:     app,
		control: control,
	}
}

func (i *injector) setService(service any, config *ServiceConfig) {
	rp := &reflectProvider{
		reflectType:  reflect.TypeOf(service),
		reflectValue: reflect.ValueOf(service),
	}
	if rp.reflectType.Kind() != reflect.Ptr {
		log.Fatalf("Service [%s] missing pointer.", rp.reflectType.Name())
	}
	i.app.services[rp.reflectValue.Elem().Type().String()] = &appService{
		provider: rp,
		config:   config,
	}
}

func (i *injector) autoinject(dst reflect.Value) {
	if dst.Kind() == reflect.Ptr {
		dst = dst.Elem()
	}
	fieldsLen := dst.NumField()
	if fieldsLen == 0 {
		return
	}
	for index := 0; index < fieldsLen; index++ {
		field := dst.Field(index)
		if field.Type().Kind() != reflect.Ptr {
			continue
		}
		switch field.Type().Elem() {
		case controllerType:
			i.injectController(field)
		default:
			i.createNewInstance(field)
		}
		if field.Elem().FieldByName(componentType.Name()).IsValid() {
			i.injectComponent(field)
		}
		if field.Elem().FieldByName(serviceType.Name()).IsValid() {
			i.injectService(field)
		}
	}
}

func (i *injector) createNewInstance(dst reflect.Value) {
	if !dst.IsNil() {
		return
	}
	dst.Set(reflect.New(dst.Type().Elem()))
}

func (i *injector) injectController(dst reflect.Value) {
	dst.Set(reflect.ValueOf(&Controller{Control: i.control}))
}

func (i *injector) injectComponent(dst reflect.Value) {
	c := &Component{
		Control: i.control,
		Handle:  make(map[string]string),
	}
	dst.Elem().FieldByName(componentType.Name()).Set(
		reflect.ValueOf(c),
	)
}

func (i *injector) injectService(dst reflect.Value) {
	name := dst.Type().Elem().String()
	if service, ok := i.app.services[name]; ok {
		if service.config.Root {
			dst.Set(service.provider.reflectValue)
		}
		if !service.config.Root {
			newService := reflect.New(service.provider.reflectType.Elem())
			i.autoinject(newService)
			dst.Set(newService)
		}
		s := &Service{
			control: i.control,
		}
		dst.Elem().FieldByName(serviceType.Name()).Set(
			reflect.ValueOf(s),
		)
	}
}

func (i *injector) injectServices(dst reflect.Value) {
	if dst.Kind() == reflect.Ptr {
		dst = dst.Elem()
	}
	fieldsLen := dst.NumField()
	if fieldsLen == 0 {
		return
	}
	for index := 0; index < fieldsLen; index++ {
		if dst.Field(index).Type().Kind() == reflect.Ptr {
			continue
		}
		name := dst.Field(index).Type().String()
		if service, ok := i.app.services[name]; ok {
			if service.config.Root {
				dst.Field(index).Set(service.provider.reflectValue.Elem())
			}
			if !service.config.Root {
				newService := reflect.New(service.provider.reflectType.Elem())
				i.injectServices(newService)
				dst.Field(index).Set(newService.Elem())
			}
		}
	}
}

func (i *injector) copy(dst reflect.Value, src reflect.Value) {
	if dst.Kind() == reflect.Ptr {
		dst = dst.Elem()
	}
	if src.Kind() == reflect.Ptr {
		src = src.Elem()
	}
	if src.Kind() != reflect.Struct && dst.Kind() != reflect.Struct {
		return
	}
	dstFieldsLen := dst.NumField()
	if dstFieldsLen == 0 {
		return
	}
	for index := 0; index < dstFieldsLen; index++ {
		if dst.Field(index).Type().Kind() == reflect.Ptr {
			switch dst.Field(index).Type().Elem() {
			case controllerType:
				cc, ok := i.firstFieldByType(src, controllerType)
				if ok && !dst.Field(index).IsNil() {
					dst.Field(index).Elem().Set(cc)
				}
			}
			continue
		}
		dst.Field(index)
		fieldName := dst.Field(index).Type().Name()
		if len(fieldName) == 0 {
			continue
		}
		if src.FieldByName(fieldName).IsValid() {
			dst.Field(index).Set(src.FieldByName(fieldName))
		}
	}
}

func (i *injector) firstFieldByType(src reflect.Value, fieldType reflect.Type) (reflect.Value, bool) {
	fieldsLen := src.NumField()
	if fieldsLen == 0 {
		return reflect.Value{}, false
	}
	for index := 0; index < fieldsLen; index++ {
		if src.Field(index).Type().Kind() != reflect.Ptr {
			continue
		}
		if src.Field(index).Type().Elem() == fieldType {
			return src.Field(index), true
		}
	}
	return reflect.Value{}, false
}
