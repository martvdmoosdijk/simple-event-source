package simpleventsrc

import (
	"reflect"
)

func nameOf(iface interface{}) string {
	if t := reflect.TypeOf(iface); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}
