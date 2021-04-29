package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var fre = regexp.MustCompile(`[A-Z]+[^A-Z]*`)

func fmtName(name string) string {
	toks := fre.FindAllString(name, -1)
	tmp := strings.Join(toks, " ")
	return strings.Title(tmp)
}

func fmtType(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()

	case reflect.Int:
		i := v.Interface().(int)
		return strconv.FormatInt(int64(i), 10)

	default:
		return "unsupported type"
	}
}

func generateMessage(o interface{}) map[string]string {
	var ret = make(map[string]string)
	var tmp = make(map[string]*strings.Builder)

	e := reflect.ValueOf(o).Elem()

	for i := 0; i < e.NumField(); i++ {
		tag := e.Type().Field(i).Tag

		section := tag.Get("section")
		if section == "" {
			section = "main"
		}

		b, ok := tmp[section]
		if !ok {
			b = &strings.Builder{}
			tmp[section] = b
		}

		name := fmtName(e.Type().Field(i).Name)
		typ := fmtType(e.Field(i))
		b.WriteString(fmt.Sprintf("%s: *%s*\n", name, typ))
	}

	for k, v := range tmp {
		ret[k] = v.String()
	}

	return ret
}
