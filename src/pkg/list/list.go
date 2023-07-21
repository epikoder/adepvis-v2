package list

import (
	"reflect"
)

func Contains(a interface{}, i interface{}) (bool) {
	val := reflect.ValueOf(a)
	if val.Kind() != reflect.Slice {
		return false
	}

	for j := 0; j < val.Len(); j++  {
		if val.Index(j).Interface() == i {
			return true
		}
	}
	return false
}