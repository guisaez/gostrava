package gostrava

import (
	"bytes"
	"fmt"
	"reflect"
)

var timeStampType = reflect.TypeOf(TimeStamp{})

// Stringify attempts to create a reasonable string representation of types in
// the GitHub library. It does things like resolve pointers to their values
// and omits struct fields with nil values.
func Stringify(message interface{}) string {
	var buf bytes.Buffer
	v := reflect.ValueOf(message)
	stringifyValue(&buf, v, 0)
	return buf.String()
}

// stringifyValue was heavily inspired by the goprotobuf library.

func stringifyValue(w *bytes.Buffer, val reflect.Value, indent int) {
	indentation := func(level int) {
		for i := 0; i < level; i++ {
			w.Write([]byte("  ")) // Using two spaces for each indentation level
		}
	}

	if val.Kind() == reflect.Ptr && val.IsNil() {
		indentation(indent)
		w.Write([]byte("<nil>"))
		return
	}

	v := reflect.Indirect(val)

	switch v.Kind() {
	case reflect.String:
		indentation(indent)
		fmt.Fprintf(w, `"%s"`, v)
	case reflect.Slice:
		indentation(indent)
		w.Write([]byte{'['})
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				w.Write([]byte{' '})
			}

			stringifyValue(w, v.Index(i), indent+1)
		}
		indentation(indent)
		w.Write([]byte{']'})
	case reflect.Struct:
		indentation(indent)
		if v.Type().Name() != "" {
			w.Write([]byte(v.Type().String()))
		}

		// special handling of Timestamp values
		if v.Type() == timeStampType {
			fmt.Fprintf(w, "{%s}", v.Interface())
			return
		}

		w.Write([]byte{'{'})
		w.Write([]byte{'\n'})

		var sep bool
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				continue
			}
			if fv.Kind() == reflect.Slice && fv.IsNil() {
				continue
			}
			if fv.Kind() == reflect.Map && fv.IsNil() {
				continue
			}

			if sep {
				w.Write([]byte(",\n"))
			} else {
				sep = true
			}

			indentation(indent + 1)
			w.Write([]byte(v.Type().Field(i).Name))
			w.Write([]byte{':'})
			stringifyValue(w, fv, indent+1)
		}
		w.Write([]byte{'\n'})
		indentation(indent)
		w.Write([]byte{'}'})
	default:
		indentation(indent)
		if v.CanInterface() {
			fmt.Fprint(w, v.Interface())
		}
	}
}
