package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

var (
	MaxDisplayLevel = 5
)

func Display(name string, x interface{}) {
	fmt.Printf("Display %s (%T):\n", name, x)
	display(name, reflect.ValueOf(x), 0)
}

func display(path string, v reflect.Value, level int) {
	if level > MaxDisplayLevel {
		fmt.Printf("%s = %s\n", path, formatAtom(v))
		return
	}
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i), level+1)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			display(fieldPath, v.Field(i), level+1)
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			display(fmt.Sprintf("%s[%s]", path, formatMapKey(key)), v.MapIndex(key), level+1)
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display(fmt.Sprintf("(*%s)", path), v.Elem(), level+1)
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			display(path+".value", v.Elem(), level+1)
		}
	default:
		fmt.Printf("%s = %s\n", path, formatAtom(v))
	}
}

func formatMapKey(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Struct:
		b := &bytes.Buffer{}
		b.WriteByte('{')
		for i := 0; i < v.NumField(); i++ {
			if i != 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(b, "%s: %s", v.Type().Field(i).Name, formatAtom(v.Field(i)))
		}
		b.WriteByte('}')
		return b.String()
	case reflect.Array:
		b := &bytes.Buffer{}
		b.WriteByte('{')
		for i := 0; i < v.Len(); i++ {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(formatAtom(v.Index(i)))
		}
		b.WriteByte('}')
		return b.String()
	default:
		return formatAtom(v)
	}
}

func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	}
	return v.Type().String() + " value"
}

type Movie struct {
	Title, Subtitle string
	Year            int
	Color           bool
	Actor           map[string][]string
	Oscars          []string
	Sequel          *string
	Tests           map[TestStruct][]string
}

type TestStruct struct {
	Foo string
}

func main() {
	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Color:    false,
		Actor: map[string][]string{
			"Dr. Strangelove":            []string{"Peter Sellers", "happy"},
			"Grp. Capt. Lionel Mandrake": []string{"Peter Sellers"},
			"Pres. Merkin Muffley":       []string{"Peter Sellers"},
			"Gen. Buck Turgidson":        []string{"George C. Scott"},
			"Brig. Gen. Jack D. Ripper":  []string{"Sterling Hayden"},
			`Maj. T.J. "King" Kong`:      []string{"Slim Pickens"},
		},
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
		Tests: map[TestStruct][]string{
			TestStruct{Foo: "foo"}: []string{"hello", "world"},
		},
	}
	Display("strangelove", strangelove)
}
