package otto

import (
	"encoding/json"
	"reflect"
)

type GoObjectHandler interface {
	GetValue(name string) interface{}
	GetProperty(name string) (interface{}, bool)
	Enumerate(all bool, each func(string) bool)
}

type GoObjectSetter interface {
	SetValue(name string, value interface{}) error
	CanSetValue(name string) bool
}

type GoObjectExtender interface {
}

type _GoHandlerObject struct {
	value GoObjectHandler
}

type GoObjectComparable interface {
	CompareWith(b interface{}) LessThanResult
}

func goHandlerIsCompareAble(self *_object) bool {
	object := self.value.(_GoHandlerObject)
	if _, ok := object.value.(GoObjectComparable); ok {
		return true
	}
	return false
}

func goHandlerCompareWith(self *_object, y *_object) LessThanResult {
	object := self.value.(_GoHandlerObject)
	if c, ok := object.value.(GoObjectComparable); ok {
		return c.CompareWith(y.value)
	}
	return LessThanUndefined
}

func goHandlerGetOwnProperty(self *_object, name string) *_property {
	object := self.value.(_GoHandlerObject)
	value := object.value.GetValue(name)
	var rv = reflect.ValueOf(value)
	if rv.IsValid() {
		return &_property{self.runtime.toValue(value), 0110}
	}

	return objectGetOwnProperty(self, name)
}

func goHandlerCanPut(self *_object, name string) bool {
	object := self.value.(_GoHandlerObject)
	if setter, ok := object.value.(GoObjectSetter); ok {
		return setter.CanSetValue(name)
	}
	return false
}

func goHandlerPut(self *_object, name string, value Value, throw bool) {
	object := self.value.(_GoHandlerObject)
	if setter, ok := object.value.(GoObjectSetter); ok {
		setter.SetValue(name, value.value)
	}
}

func goHandlerEnumerate(self *_object, all bool, each func(string) bool) {
	object := self.value.(_GoHandlerObject)
	object.value.Enumerate(all, each)

	objectEnumerate(self, all, each)
}

func goHandlerMarshalJSON(self *_object) json.Marshaler {
	object := self.value.(*_GoHandlerObject)
	switch marshaler := object.value.(type) {
	case json.Marshaler:
		return marshaler
	}
	return nil
}
