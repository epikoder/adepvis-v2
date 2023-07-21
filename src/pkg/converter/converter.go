package converter

import (
	"fmt"
	"reflect"
	"strconv"
)

func GetInt(i interface{}) (int, error) {
	v := reflect.ValueOf(i)
	s := fmt.Sprintf("%v", v)
	return strconv.Atoi(s)
}
