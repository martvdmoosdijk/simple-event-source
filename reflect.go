package simpleventsrc

import (
	"reflect"
)

func typeOf(item interface{}) reflect.Type {
	if t := reflect.TypeOf(item); t.Kind() == reflect.Ptr {
		return t.Elem()
	} else {
		return t
	}
}
