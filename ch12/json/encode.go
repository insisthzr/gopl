package json

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

func encode(buf *bytes.Buffer, v reflect.Value, indent int) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("null")

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(buf, "%d", v.Uint())

	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%g", v.Float())

	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())

	case reflect.Ptr:
		return encode(buf, v.Elem(), indent)

	case reflect.Bool:
		if v.Bool() {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}

	case reflect.Array, reflect.Slice:
		buf.WriteByte('[')
		for i := 0; i < v.Len(); i++ {
			if i != 0 {
				buf.WriteByte(',')
			}
			err := encode(buf, v.Index(i), indent)
			if err != nil {
				return err
			}
		}
		buf.WriteByte(']')

	case reflect.Struct:
		indent += 2
		buf.WriteString("{\n")
		for i := 0; i < v.NumField(); i++ {
			buf.WriteString(strings.Repeat(" ", indent))
			fmt.Fprintf(buf, "%q:", v.Type().Field(i).Name)
			err := encode(buf, v.Field(i), indent)
			if err != nil {
				return err
			}
			if i != v.NumField()-1 {
				buf.WriteString(",")
			}
			buf.WriteString("\n")
		}
		buf.WriteString("}\n")

	case reflect.Map:
		indent += 2
		buf.WriteString("{\n")
		for i, key := range v.MapKeys() {
			buf.WriteString(strings.Repeat(" ", indent))
			err := encode(buf, key, indent)
			if err != nil {
				return err
			}
			buf.WriteByte(':')
			err = encode(buf, v.MapIndex(key), indent)
			if err != nil {
				return err
			}
			if i != len(v.MapKeys())-1 {
				buf.WriteString(",")
			}
			buf.WriteString("\n")
		}
		buf.WriteString("}\n")

	default:
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := encode(&buf, reflect.ValueOf(v), 0)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
