package utils

import (
	"reflect"
)

func isASCIIUpper(r rune) bool {
	return 'A' <= r && r <= 'Z'
}

func ConvertToSlice(input interface{}) (data []interface{}) {
	val := reflect.ValueOf(input)
	data = make([]interface{}, 0)
start:
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		data = make([]interface{}, val.Len())
		for i := 0; i < val.Len(); i++ {
			data = append(data, val.Index(i).Interface())
		}
	case reflect.Ptr:
		val = val.Elem()
		goto start
	}
	return
}
