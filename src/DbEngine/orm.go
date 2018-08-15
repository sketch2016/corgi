package DbEngine

import (
	"reflect"
	"strings"
)

const (
	IN int = iota
	IN_OUT
	OUT
)

type ormValue struct {
	value *reflect.Value
	inout int
}

//GetInstance instance func
func ormReflect(instance interface{}) (r1 map[string]*ormValue, r2 reflect.Value) {
	r1 = make(map[string]*ormValue)

	reflectTag(instance, &r1)
	r2 = reflect.ValueOf(instance).MethodByName("ConvertValue")
	return
}

func reflectTag(instance interface{}, r *map[string]*ormValue) {
	v := reflect.ValueOf(instance)
	f := reflect.TypeOf(instance)

	if f.Kind() == reflect.Ptr {
		v = v.Elem()
		f = f.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i)
		if value.Kind() == reflect.Struct {
			reflectTag(value, r)
			continue
		} else {
			field := f.Field(i)
			if len(field.Tag) != 0 {
				rr := *r
				vret := new(ormValue)
				result1, result2 := parseTag(string(field.Tag))
				vret.inout = result2
				vret.value = &value
				rr[result1] = vret
			}
		}
	}

	//fmt.Println("rr is ", r)
}

func parseTag(tag string) (string, int) {
	str := strings.Split(tag, ":")
	var inout int
	if len(str) > 1 {
		switch str[1] {
		case "in":
			inout = IN

		case "out":
			inout = OUT

		case "inout":
			inout = IN_OUT
		}
	} else {
		inout = IN_OUT
	}

	return str[0], inout
}
