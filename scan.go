package haproxy

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	truth_list = []string{"true", "1", "yes", "y", "on"}
	false_list = []string{"false", "0", "no", "n", "off"}
	bool_map   map[string]bool
)

func init() {
	bool_map = make(map[string]bool, 0)
	for _, k := range truth_list {
		bool_map[k] = true
	}
	for _, k := range false_list {
		bool_map[k] = false
	}
}
func BoolVal(s string) (bool, error) {
	k := strings.ToLower(s)
	v, found := bool_map[k]
	if found {
		return v, nil
	}
	return v, fmt.Errorf("invalid value %s", s)
}

func ScanMap(m map[string]string, dest interface{}) error {
	v := reflect.ValueOf(dest)
	base := reflect.Indirect(v)
	for i := 0; i < base.NumField(); i++ {
		df := base.Field(i)
		dt := df.Type()
		var field_name string
		tag := base.Type().Field(i).Tag.Get("scan")
		if tag == "" {
			// field_name = strings.ToLower(base.Type().Field(i).Name)
			field_name = base.Type().Field(i).Name
		} else {
			field_name = tag
		}
		v, ok := m[field_name]
		if ok == false {
			continue
		}

		if v == "" {
			continue
		}

		switch dt.Kind() {

		case reflect.Int:
			nv, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return fmt.Errorf("%s parameter: parse error", field_name)
			}
			df.SetInt(nv)

		case reflect.Bool:
			b, err := BoolVal(v)
			if err != nil {
				return fmt.Errorf("%s parameter: parse bool error", field_name)
			}
			df.SetBool(b)

		case reflect.String:
			df.SetString(v)

		}

	}
	return nil
}
