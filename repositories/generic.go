package repositories

import "reflect"

func MapToEntCreate[T any](builder any, data T) any {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i).Name
		value := v.Field(i).Interface()

		method := reflect.ValueOf(builder).MethodByName("Set" + field)
		if method.IsValid() {
			method.Call([]reflect.Value{reflect.ValueOf(value)})
		}
	}

	return builder
}
