package util

import "reflect"

func ConvertToStruct(m map[string]interface{}, s interface{}) {
    v := reflect.Indirect(reflect.ValueOf(s))

    for i := 0; i < v.NumField(); i++ {
        key := v.Type().Field(i).Name
        v.Field(i).Set(reflect.ValueOf(m[key]))
    }
}

func ConvertToMap(s interface{}) map[string]interface{} {
    m := make(map[string]interface{})

    v := reflect.ValueOf(s)

    for i := 0; i < v.NumField(); i++ {
        key := v.Type().Field(i).Name
        val := v.Field(i).Interface()

        m[key] = val
    }

    return m
}

func StructName(s interface{}) string {
    v := reflect.TypeOf(s)
    for ;v.Kind() == reflect.Ptr; {
        v = v.Elem()
    }
    return v.Name()
}
