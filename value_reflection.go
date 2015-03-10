package cdnlysis

import (
	"errors"
	"log"
	"reflect"
	"strconv"
)

func throw(err error) {
	panic(err)
}

type conversionError struct {
	Value string
	Type  reflect.Type
}

func (e conversionError) Error() {
	throw(errors.New("json: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()))
}

type Number string

// String returns the literal text of the number.
func (n Number) String() string { return string(n) }

// Float64 returns the number as a float64.
func (n Number) Float64() (float64, error) {
	return strconv.ParseFloat(string(n), 64)
}

// Int64 returns the number as an int64.
func (n Number) Int64() (int64, error) {
	return strconv.ParseInt(string(n), 10, 64)
}

var numberType = reflect.TypeOf(Number(""))

func convertValue(v reflect.Value, actualValue string) {
	switch v.Kind() {
	default:
		log.Println(v.Kind(), v.Type())
		if v.Kind() == reflect.String && v.Type() == numberType {
			v.SetString(actualValue)
			break
		}

		conversionError{"number", v.Type()}.Error()
		break

	case reflect.String:
		v.SetString(string(actualValue))
		break

	case reflect.Interface:
		n, err := strconv.ParseFloat(actualValue, 64)
		if err != nil {
			throw(err)
			break
		}

		if v.NumMethod() != 0 {
			conversionError{"number", v.Type()}.Error()
			break
		}

		v.Set(reflect.ValueOf(n))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(actualValue, 10, 64)
		if err != nil || v.OverflowInt(n) {
			conversionError{"number " + actualValue, v.Type()}.Error()
			break
		}

		v.SetInt(n)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(actualValue, 10, 64)
		if err != nil || v.OverflowUint(n) {
			conversionError{"number " + actualValue, v.Type()}.Error()
			break
		}
		v.SetUint(n)

	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(actualValue, v.Type().Bits())
		if err != nil || v.OverflowFloat(n) {
			conversionError{"number " + actualValue, v.Type()}.Error()
			break
		}
		v.SetFloat(n)
	}

}
