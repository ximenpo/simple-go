package simple

import (
	"errors"
	"reflect"
)

func Infect(struct_obj interface{}, value interface{}) (int, error) {
	if err := _infectCheckParams(struct_obj, value); err != nil {
		return 0, err
	}

	V := reflect.ValueOf(struct_obj).Elem()
	return _infectStructFields(
		&V,
		"",
		reflect.TypeOf(value).Elem(),
		reflect.ValueOf(value).Elem(),
	)
}

func InfectByName(struct_obj interface{}, value interface{}, name string) (int, error) {
	if err := _infectCheckParams(struct_obj, value); err != nil {
		return 0, err
	}

	V := reflect.ValueOf(struct_obj).Elem()
	return _infectStructFields(
		&V,
		name,
		reflect.TypeOf(value).Elem(),
		reflect.ValueOf(value).Elem(),
	)
}

func _infectCheckObj(struct_obj interface{}) error {
	if struct_obj == nil {
		return errors.New("struct_obj must not be nil")
	}
	T := reflect.TypeOf(struct_obj)
	if T.Kind() != reflect.Ptr {
		return errors.New("struct_obj must be pointer type")
	}
	if T.Elem().Kind() != reflect.Struct {
		return errors.New("struct_obj must be struct object pointer")
	}
	return nil
}

func _infectCheckParams(struct_obj interface{}, value interface{}) error {
	if err := _infectCheckObj(struct_obj); err != nil {
		return nil
	}

	if value == nil {
		return errors.New("value must not be nil")
	}
	if reflect.TypeOf(value).Kind() != reflect.Ptr {
		return errors.New("value must be pointer type")
	}

	return nil
}

func _infectStructFields(obj *reflect.Value, name string, T reflect.Type, V reflect.Value) (int, error) {
	child_sum := obj.NumField()
	infect_sum := 0
	check_name := (len(name) > 0)
	OT := obj.Type()
	for i := 0; i < int(child_sum); i++ {
		field := obj.Field(i)
		if !field.CanSet() {
			continue
		}

		// field assign
		if field.Type() == T && (!check_name || name == OT.Field(i).Name) {
			field.Set(V)
			infect_sum++
			continue
		}

		// deep field assign
		var pfv *reflect.Value
		switch field.Kind() {
		case reflect.Struct:
			pfv = &field
		case reflect.Interface:
			if !field.IsNil() {
				fv := reflect.ValueOf(field.Interface()).Elem()
				pfv = &fv
			}
		case reflect.Ptr:
			if !field.IsNil() && field.Elem().Kind() == reflect.Struct {
				fv := field.Elem()
				pfv = &fv
			}
		}
		if pfv == nil {
			continue
		}
		if sum, err := _infectStructFields(pfv, name, T, V); err != nil {
			return infect_sum + sum, err
		} else {
			infect_sum += sum
		}
	}

	return infect_sum, nil
}
