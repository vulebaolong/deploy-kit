package ui

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

func Step(message string) {
	fmt.Printf("\n%s%s▶ STEP%s %s\n", Bold, Cyan, Reset, message)
}

func Success(message string) {
	fmt.Printf("%s✔ OK%s %s\n", Green, Reset, message)
}

func Warn(message string) {
	fmt.Printf("%s⚠ WARN%s %s\n", Yellow, Reset, message)
}

func Error(message string) {
	fmt.Printf("%s✖ ERROR%s %s\n", Red, Reset, message)
}

func Info(message string) {
	fmt.Printf("%s[INFO]%s %s\n", Blue, Reset, message)
}

// ui.Step(fmt.Sprintf("Remove old image: %s", imageFullName))
// ui.Info(fmt.Sprintf("Dockerfile: %s", dockerfilePath))
// ui.Success("Build image thành công")
// ui.Warn("Image cũ không tồn tại, bỏ qua")
// ui.Error("Build image thất bại")

func PrintStruct(title string, data any) {
	fmt.Printf("\n====== %s ======\n", title)
	printValue(reflect.ValueOf(data), 0)
	fmt.Println("==============================")
}

func printValue(v reflect.Value, level int) {
	if !v.IsValid() {
		return
	}

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		fmt.Printf("%s%v\n", indent(level), v.Interface())
		return
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanInterface() {
			continue
		}

		fieldName := fieldType.Name

		switch field.Kind() {
		case reflect.Struct:
			fmt.Printf("%s%s\n", indent(level), fieldName)
			printValue(field, level+1)

		case reflect.Ptr:
			if field.IsNil() {
				fmt.Printf("%s%s: <nil>\n", indent(level), fieldName)
			} else if field.Elem().Kind() == reflect.Struct {
				fmt.Printf("%s%s\n", indent(level), fieldName)
				printValue(field.Elem(), level+1)
			} else {
				fmt.Printf("%s%s: %v\n", indent(level), fieldName, field.Elem().Interface())
			}

		case reflect.Map:
			if field.IsNil() || field.Len() == 0 {
				fmt.Printf("%s%s: <empty>\n", indent(level), fieldName)
				continue
			}

			fmt.Printf("%s%s\n", indent(level), fieldName)

			for _, key := range field.MapKeys() {
				value := field.MapIndex(key)
				fmt.Printf("%s%v: %v\n", indent(level+1), key.Interface(), value.Interface())
			}

		case reflect.Slice, reflect.Array:
			if field.Len() == 0 {
				fmt.Printf("%s%s: <empty>\n", indent(level), fieldName)
				continue
			}

			fmt.Printf("%s%s\n", indent(level), fieldName)
			for j := 0; j < field.Len(); j++ {
				item := field.Index(j)
				if item.Kind() == reflect.Struct || (item.Kind() == reflect.Ptr && !item.IsNil()) {
					printValue(item, level+1)
				} else {
					fmt.Printf("%s- %v\n", indent(level+1), item.Interface())
				}
			}

		default:
			fmt.Printf("%s%s: %v\n", indent(level), fieldName, field.Interface())
		}
	}
}

func indent(level int) string {
	return strings.Repeat("  ", level)
}
