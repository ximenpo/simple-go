package simple

import (
	"errors"
    "log"
	"reflect"
)

func Infect_InterfaceField(struct_obj interface{}, value interface{}) error {
	if struct_obj == nil {
		return errors.New("struct_obj must not be nil")
	}
	T := reflect.TypeOf(struct_obj)
	if T.Kind() != reflect.Ptr {
		return errors.New("struct_obj must be pointer type")
	}
	V := reflect.ValueOf(struct_obj).Elem()
	if V.Kind() != reflect.Struct {
		return errors.New("struct_obj must be struct object pointer")
	}

	XT := reflect.TypeOf(value)
	XV := reflect.ValueOf(value)

	return _infectInterfaceField(&V, XT, XV)
}

func _infectInterfaceField(obj *reflect.Value, T reflect.Type, V reflect.Value) error {
log.Println("processing type ", T.Kind(), T, "for", obj.Kind(), obj.Type())
	child_sum := obj.NumField()
	for i := 0; i < int(child_sum); i++ {
		field := obj.Field(i)
        if !field.CanSet(){
            continue
        }
        log.Println("field ", field.Kind(), field.Type())
		switch field.Kind() {
        case reflect.Ptr:
			if field.Type() == T {
				field.Set(V)
                log.Println("!setted", field.Kind())
			}else if !field.IsNil() && field.Type().Elem().Kind() == reflect.Struct {
    			if err := _infectInterfaceField(&field, T, V); err != nil {
    				return err
    			}
            }
        case reflect.Interface:
            if V.CanInterface() && T.Implements(field.Type()) {
                field.Set(V)
                log.Println("!setted", field.Kind())
            }else if !field.IsNil() {
                SV := reflect.ValueOf(field.Interface()).Elem()
    			if err := _infectInterfaceField(&SV, T, V); err != nil {
    				return err
    			}
            }
		case reflect.Struct:
			if err := _infectInterfaceField(&field, T, V); err != nil {
				return err
			}
        default:
		}
	}

	return nil
}
