package util

import (
	"errors"
	"reflect"
	"strings"
)

// ConvertCopyStruct copies the values of the fields in src matching values in dst
// a value is copied if the field in src and dst have the same name (regardless of capitalization)
// and if they have same type, or their types can be converted to one another (e.g. int32 and int64).
// or if a ConverterFunc is supplied to map between the src and dst field
// Both src and dst should be pointers to structs.
//
// ConvertCopyStruct should be primarily used to conveniently convert structs with similar fields,
// e.g. database models generated with sqlboiler and rpc models generated with protoc-gen-gogo
//
// ConvertCopyStruct currently does not handle nested structs, or pointer fields (this includes maps, channels, slices, etc.)
// and is basically a worse version of github.com/jinzhu/copier (which unfortunately does not ignore capitalization,
// which is the only justification for writing this janky ass function)
func ConvertCopyStruct(dst, src interface{}, converters map[string]ConverterFunc) error {
	srcT := reflect.TypeOf(src)
	dstT := reflect.TypeOf(dst)

	if srcT.Kind() != reflect.Ptr ||
		dstT.Kind() != reflect.Ptr ||
		srcT.Elem().Kind() != reflect.Struct ||
		dstT.Elem().Kind() != reflect.Struct {

		return errors.New("parameters to CopyConvertStruct should be pointers to structs")
	}

	srcV := reflect.ValueOf(src).Elem()
	dstV := reflect.ValueOf(dst).Elem()

	srcFieldMetas := []reflect.StructField{}
	dstFieldMetas := []reflect.StructField{}

	for i := 0; i < srcT.Elem().NumField(); i++ {
		srcFieldMeta := srcT.Elem().Field(i)
		srcFieldMetas = append(srcFieldMetas, srcFieldMeta)
	}

	for i := 0; i < dstT.Elem().NumField(); i++ {
		dstFieldMeta := dstT.Elem().Field(i)
		dstFieldMetas = append(dstFieldMetas, dstFieldMeta)
	}

	for _, srcFieldMeta := range srcFieldMetas {
		// get the index of a field in dstfields that has the same (case-insensitive) name
		// and the src field value can be assigned to the dst field
		matchingFieldIdx := -1
		for i, dstFieldMeta := range dstFieldMetas {
			if strings.EqualFold(srcFieldMeta.Name, dstFieldMeta.Name) {
				matchingFieldIdx = i
			}
		}

		// ignore if there is no match
		if matchingFieldIdx == -1 {
			continue
		}

		dstFieldMeta := dstFieldMetas[matchingFieldIdx]

		// check if the field is directly assignable or a converter is present
		var converter ConverterFunc

		for k, v := range converters {
			if strings.EqualFold(k, dstFieldMeta.Name) {
				converter = v
				break
			}
		}

		if !dstFieldMeta.Type.AssignableTo(srcFieldMeta.Type) && converter == nil {
			continue
		}

		dstField := dstV.FieldByName(dstFieldMeta.Name)
		srcField := srcV.FieldByName(srcFieldMeta.Name)

		if !dstField.CanSet() {
			continue
		}

		if converter != nil {
			dstField.Set(
				reflect.ValueOf(converter(srcField.Interface())),
			)
			continue
		}

		dstField.Set(srcField)
	}

	return nil
}
