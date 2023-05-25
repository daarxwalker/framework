package framework

import "reflect"

type reflectProvider struct {
	reflectType  reflect.Type
	reflectValue reflect.Value
}
