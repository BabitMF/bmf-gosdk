package builder

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func combineStreamsWithoutNil(streams ...*BMFStream) []*BMFStream {
	var ret []*BMFStream
	for _, s := range streams {
		if s != nil {
			ret = append(ret, s)
		}
	}
	return ret
}

func combineStreamWithNil(streams ...*BMFStream) []*BMFStream {
	return streams
}

func dumpFilterOption(option interface{}) string {
	return dumpFilter(reflect.ValueOf(option))
}

func dumpFilter(v reflect.Value) string {
	var res []string
	switch v.Kind() {
	case reflect.Interface:
		return dumpFilter(v.Elem())
	case reflect.Ptr:
		return dumpFilter(v.Elem())
	case reflect.Struct:
		typ := v.Type()
		fc := typ.NumField()
		for i := 0; i < fc; i++ {
			f := typ.Field(i)
			tag := f.Tag.Get("bmf")
			if len(tag) != 0 {
				res = append(res, fmt.Sprintf("%s=%s", tag, dumpFilter(v.Field(i))))
			} else if f.PkgPath == "" || f.Anonymous {
				res = append(res, fmt.Sprintf("%s=%s", f.Name, dumpFilter(v.Field(i))))
			} else {
				res = append(res, dumpFilter(v.Field(i)))
			}
		}
	case reflect.Map:
		keys := v.MapKeys()
		for _, k := range keys {
			res = append(res, fmt.Sprintf("%s=%s", dumpFilter(k), dumpFilter(v.MapIndex(k))))
		}
	case reflect.Slice, reflect.Array:
		l := v.Len()
		for i := 0; i < l; i++ {
			res = append(res, dumpFilter(v.Index(i)))
		}
	case reflect.String:
		res = append(res, v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		res = append(res, strconv.FormatInt(v.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		res = append(res, strconv.FormatUint(v.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		precision := 64
		if v.Kind() == reflect.Float32 {
			precision = 32
		}
		s := strconv.FormatFloat(v.Float(), 'g', -1, precision)
		switch s {
		case "+Inf":
			s = ".inf"
		case "-Inf":
			s = "-.inf"
		case "NaN":
			s = ".nan"
		}
		res = append(res, s)
	case reflect.Bool:
		if v.Bool() {
			res = append(res, "true")
		} else {
			res = append(res, "false")
		}
	default:
		panic("Unknown type, cannot gen para for module, Type =" + v.Kind().String())
	}
	return strings.Join(res, ":")
}

func data2Map(inData interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	if inData == nil {
		return res
	}

	inV := reflect.ValueOf(inData)
	switch inV.Kind() {
	case reflect.Map:
		keys := inV.MapKeys()
		for _, k := range keys {
			if k.Kind() != reflect.String {
				panic("option map key should be string!")
			}
			res[k.String()] = inV.MapIndex(k).Interface()
		}
	case reflect.Struct:
		typ := inV.Type()
		fc := typ.NumField()
		for i := 0; i < fc; i++ {
			f := typ.Field(i)
			bmfTag := f.Tag.Get("bmf")
			jsonTag := f.Tag.Get("json")
			if len(bmfTag) != 0 {
				res[bmfTag] = inV.Field(i).Interface()
			} else if len(jsonTag) != 0 {
				res[jsonTag] = inV.Field(i).Interface()
			} else if f.PkgPath == "" || f.Anonymous {
				res[f.Name] = inV.Field(i).Interface()
			}
		}
	case reflect.Ptr:
		return data2Map(inV.Elem().Interface())
	default:
		panic("option type only supports map or struct")
	}
	return res
}
