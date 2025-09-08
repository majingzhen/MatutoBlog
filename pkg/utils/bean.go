package utils

import (
	"fmt"
	"reflect"
)

// CopyProperties 拷贝源对象的属性到目标对象
// src: 源对象(指针或值)
// dst: 目标对象(必须是指针)
func CopyProperties(src, dst interface{}) error {
	if src == nil || dst == nil {
		return fmt.Errorf("源对象或目标对象不能为nil")
	}

	srcValue := reflect.ValueOf(src)
	dstValue := reflect.ValueOf(dst)

	// 目标必须是指针
	if dstValue.Kind() != reflect.Ptr {
		return fmt.Errorf("目标对象必须是指针")
	}

	// 获取指针指向的值
	dstElem := dstValue.Elem()
	if !dstElem.CanSet() {
		return fmt.Errorf("目标对象不可设置")
	}

	// 如果源对象是指针，获取其指向的值
	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}

	return copyStruct(srcValue, dstElem)
}

// CopySlice 拷贝源切片到目标切片
// src: 源切片
// dst: 目标切片(必须是指针)
func CopySlice(src, dst interface{}) error {
	if src == nil || dst == nil {
		return fmt.Errorf("源切片或目标切片不能为nil")
	}

	srcValue := reflect.ValueOf(src)
	dstValue := reflect.ValueOf(dst)

	// 目标必须是指针
	if dstValue.Kind() != reflect.Ptr {
		return fmt.Errorf("目标切片必须是指针")
	}

	dstElem := dstValue.Elem()

	// 源和目标都必须是切片
	if srcValue.Kind() != reflect.Slice {
		return fmt.Errorf("源对象必须是切片")
	}
	if dstElem.Kind() != reflect.Slice {
		return fmt.Errorf("目标对象必须是切片")
	}

	srcLen := srcValue.Len()
	dstElemType := dstElem.Type().Elem()

	// 创建目标切片
	dstSlice := reflect.MakeSlice(dstElem.Type(), srcLen, srcLen)

	for i := 0; i < srcLen; i++ {
		srcItem := srcValue.Index(i)
		dstItem := reflect.New(dstElemType).Elem()

		// 如果源元素是指针，获取其指向的值
		if srcItem.Kind() == reflect.Ptr && !srcItem.IsNil() {
			srcItem = srcItem.Elem()
		}

		if err := copyStruct(srcItem, dstItem); err != nil {
			return fmt.Errorf("拷贝第%d个元素失败: %v", i, err)
		}

		dstSlice.Index(i).Set(dstItem)
	}

	dstElem.Set(dstSlice)
	return nil
}

// copyStruct 拷贝结构体字段
func copyStruct(src, dst reflect.Value) error {
	srcType := src.Type()
	dstType := dst.Type()

	if srcType.Kind() != reflect.Struct || dstType.Kind() != reflect.Struct {
		return fmt.Errorf("源和目标都必须是结构体")
	}

	// 遍历目标结构体的字段
	for i := 0; i < dstType.NumField(); i++ {
		dstField := dstType.Field(i)
		dstValue := dst.Field(i)

		// 跳过无法设置的字段
		if !dstValue.CanSet() {
			continue
		}

		// 查找源结构体中的对应字段
		srcValue, found := findField(src, dstField.Name)
		if !found {
			continue
		}

		// 拷贝字段值
		if err := copyValue(srcValue, dstValue); err != nil {
			return fmt.Errorf("拷贝字段 %s 失败: %v", dstField.Name, err)
		}
	}

	return nil
}

// findField 在结构体中查找字段
func findField(structValue reflect.Value, fieldName string) (reflect.Value, bool) {
	structType := structValue.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.Name == fieldName {
			return structValue.Field(i), true
		}
	}

	return reflect.Value{}, false
}

// copyValue 拷贝值
func copyValue(src, dst reflect.Value) error {
	if !src.IsValid() {
		return nil
	}

	srcType := src.Type()
	dstType := dst.Type()

	// 类型完全相同，直接设置
	if srcType == dstType {
		dst.Set(src)
		return nil
	}

	// 处理指针类型
	if srcType.Kind() == reflect.Ptr && dstType.Kind() == reflect.Ptr {
		if src.IsNil() {
			return nil
		}
		if dst.IsNil() {
			dst.Set(reflect.New(dstType.Elem()))
		}
		return copyValue(src.Elem(), dst.Elem())
	}

	// 源是指针，目标不是
	if srcType.Kind() == reflect.Ptr && dstType.Kind() != reflect.Ptr {
		if src.IsNil() {
			return nil
		}
		return copyValue(src.Elem(), dst)
	}

	// 源不是指针，目标是
	if srcType.Kind() != reflect.Ptr && dstType.Kind() == reflect.Ptr {
		if dst.IsNil() {
			dst.Set(reflect.New(dstType.Elem()))
		}
		return copyValue(src, dst.Elem())
	}

	// 类型可转换
	if srcType.ConvertibleTo(dstType) {
		dst.Set(src.Convert(dstType))
		return nil
	}

	// 都是结构体，递归拷贝
	if srcType.Kind() == reflect.Struct && dstType.Kind() == reflect.Struct {
		return copyStruct(src, dst)
	}

	return fmt.Errorf("无法从 %s 类型拷贝到 %s 类型", srcType, dstType)
}

// ConvertTo 将源对象转换为目标类型的新对象
func ConvertTo[T any](src interface{}) (*T, error) {
	var dst T
	err := CopyProperties(src, &dst)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

// ConvertSliceTo 将源切片转换为目标类型的新切片
func ConvertSliceTo[T any](src interface{}) ([]T, error) {
	var dst []T
	err := CopySlice(src, &dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}
