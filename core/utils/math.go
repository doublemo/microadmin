package utils

import (
	"reflect"
)

// HasElem 检查切片中是否存在某个值
func HasElem(s, elem interface{}) bool {
	arrV := reflect.ValueOf(s)
	if arrV.Kind() == reflect.Slice || arrV.Kind() == reflect.Array {
		for i := 0; i < arrV.Len(); i++ {
			if reflect.DeepEqual(elem, arrV.Index(i).Interface()) == true {
				return true
			}
		}
	}
	return false
}
