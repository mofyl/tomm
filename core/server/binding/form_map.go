package binding

import (
	"github.com/pkg/errors"
	"reflect"
	"strconv"
)

func mapForm(ptr interface{}, form map[string][]string) error {
	sInfo := sCache.get(reflect.TypeOf(ptr))
	val := reflect.ValueOf(ptr).Elem()
	var err error
	for i, fd := range sInfo.fields {
		structVal := val.Field(i)
		if !structVal.CanSet() {
			continue
		}

		if fd.name == "" {
			continue
		}
		formV, ok := form[fd.name]
		if !ok {
			if fd.hasDefault {
				structVal.Set(fd.defaultValue)
				continue
			}
		}
		// 检查Op
		if !fd.hasDefault {
			err = CheckOptions(fd.options, formV)
			if err != nil {
				return err
			}
		}

		if formV[0] == "" && fd.hasDefault {
			structVal.Set(fd.defaultValue)
			continue
		}
		if err := setWithProperType(fd.tp.Type.Kind(), formV, structVal); err != nil {
			return err
		}
	}
	return nil
}

func setWithProperType(kind reflect.Kind, value []string, dv reflect.Value) error {

	switch kind {
	case reflect.Int:
		return setIntValue(value[0], 0, dv)
	case reflect.Int8:
		return setIntValue(value[0], 8, dv)
	case reflect.Int16:
		return setIntValue(value[0], 16, dv)
	case reflect.Int32:
		return setIntValue(value[0], 32, dv)
	case reflect.Int64:
		return setIntValue(value[0], 64, dv)
	case reflect.Uint:
		return setUintValue(value[0], 0, dv)
	case reflect.Uint8:
		return setUintValue(value[0], 8, dv)
	case reflect.Uint16:
		return setUintValue(value[0], 16, dv)
	case reflect.Uint32:
		return setUintValue(value[0], 32, dv)
	case reflect.Uint64:
		return setUintValue(value[0], 64, dv)
	case reflect.Float32:
		return setFloatValue(value[0], 32, dv)
	case reflect.Float64:
		return setFloatValue(value[0], 64, dv)
	case reflect.Bool:
		return setBoolValue(value[0], dv)
	default:
		dv.SetString(value[0])
	}
	return nil
}

func setBoolValue(value string, field reflect.Value) error {
	if value == "" {
		value = "false"
	}

	boolVal, err := strconv.ParseBool(value)

	if err != nil {
		return err
	}

	field.SetBool(boolVal)
	return nil
}

func setFloatValue(value string, bit int, field reflect.Value) error {
	if value == "" {
		value = "0.0"
	}

	fValue, err := strconv.ParseFloat(value, bit)

	if err != nil {
		return nil
	}
	field.SetFloat(fValue)
	return nil
}

func setIntValue(value string, bit int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	intV, err := strconv.ParseInt(value, 10, bit)

	if err != nil {
		return errors.WithStack(err)
	}

	field.SetInt(intV)
	return nil
}

func setUintValue(value string, bit int, dv reflect.Value) error {
	if value == "" {
		value = "0"
	}

	uintV, err := strconv.ParseUint(value, 10, bit)
	if err != nil {
		return err
	}

	dv.SetUint(uintV)
	return nil
}
