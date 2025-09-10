package utils

import (
	"github.com/jinzhu/copier"
)

// CopyProperties 拷贝源对象的属性到目标对象
// src: 源对象(指针或值)
// dst: 目标对象(必须是指针)
func CopyProperties(src, dst interface{}) error {
	return copier.Copy(dst, src)
}

// CopySlice 拷贝源切片到目标切片
// src: 源切片
// dst: 目标切片(必须是指针)
func CopySlice(src, dst interface{}) error {
	return copier.Copy(dst, src)
}

// ConvertTo 将源对象转换为目标类型的新对象
func ConvertTo[T any](src interface{}) (*T, error) {
	var dst T
	err := copier.Copy(&dst, src)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

// ConvertSliceTo 将源切片转换为目标类型的新切片
func ConvertSliceTo[T any](src interface{}) ([]T, error) {
	var dst []T
	err := copier.Copy(&dst, src)
	if err != nil {
		return nil, err
	}
	return dst, nil
}
