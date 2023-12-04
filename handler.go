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

type goHandlerObject struct {
	value GoObjectHandler
}

type GoObjectComparable interface {
	CompareWith(b interface{}) lessThanResult
}

func goHandlerIsCompareAble(self *object) bool {
	object := self.value.(goHandlerObject)
	if _, ok := object.value.(GoObjectComparable); ok {
		return true
	}
	return false
}

func goHandlerCompareWith(self *object, y Value) lessThanResult {
	object := self.value.(goHandlerObject)
	if c, ok := object.value.(GoObjectComparable); ok {
		if y.object() != nil {
			return c.CompareWith(y.object().value)
		}
		return c.CompareWith(y.value)
	}
	return lessThanUndefined
}

func goHandlerGetOwnProperty(self *object, name string) *property {
	object := self.value.(goHandlerObject)
	value := object.value.GetValue(name)
	var rv = reflect.ValueOf(value)
	if rv.IsValid() {
		return &property{self.runtime.toValue(value), 0110}
	}

	return objectGetOwnProperty(self, name)
}

func goHandlerCanPut(self *object, name string) bool {
	object := self.value.(goHandlerObject)
	if setter, ok := object.value.(GoObjectSetter); ok {
		return setter.CanSetValue(name)
	}
	return false
}

func goHandlerPut(self *object, name string, value Value, throw bool) {
	object := self.value.(goHandlerObject)
	if setter, ok := object.value.(GoObjectSetter); ok {
		setter.SetValue(name, value.value)
	}
}

func goHandlerEnumerate(self *object, all bool, each func(string) bool) {
	object := self.value.(goHandlerObject)
	object.value.Enumerate(all, each)

	objectEnumerate(self, all, each)
}

func goHandlerMarshalJSON(self *object) json.Marshaler {
	object := self.value.(*goHandlerObject)
	switch marshaler := object.value.(type) {
	case json.Marshaler:
		return marshaler
	}
	return nil
}
