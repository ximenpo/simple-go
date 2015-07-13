package databuf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
)

func min(l uint64, r uint64) uint64 {
	if l > r {
		return r
	} else {
		return l
	}
}

type DataReader struct {
	Buf io.Reader
}

func NewDataReader(buf io.Reader) *DataReader {
	return &DataReader{buf}
}

func (r *DataReader) ReadTagData(tag *DataTag) (err error) {
	var data uint8
	if err = binary.Read(r.Buf, binary.BigEndian, &data); err == nil {
		tag.UnPack(data)
	}
	return
}

func (r *DataReader) ReadUintData(size_tag uint, value *uint64) (err error) {
	switch size_tag {
	case TAG_0:
		*value = 0
	case TAG_1:
		var data uint8
		if err = binary.Read(r.Buf, binary.BigEndian, &data); err == nil {
			*value = uint64(data)
		}
	case TAG_2:
		var data uint16
		if err = binary.Read(r.Buf, binary.BigEndian, &data); err == nil {
			*value = uint64(data)
		}
	case TAG_4:
		var data uint32
		if err = binary.Read(r.Buf, binary.BigEndian, &data); err == nil {
			*value = uint64(data)
		}
	default:
		var data uint64
		if err = binary.Read(r.Buf, binary.BigEndian, &data); err == nil {
			*value = uint64(data)
		}
	}
	return
}

func (r *DataReader) ReadIntData(size_tag uint, value *int64) (err error) {
	switch size_tag {
	case TAG_0:
		*value = 0
	case TAG_1:
		var data int8
		if err = binary.Read(r.Buf, binary.BigEndian, &data); err == nil {
			*value = int64(data)
		}
	case TAG_2:
		var data int16
		if err = binary.Read(r.Buf, binary.BigEndian, &data); err == nil {
			*value = int64(data)
		}
	case TAG_4:
		var data int32
		if err = binary.Read(r.Buf, binary.BigEndian, &data); err == nil {
			*value = int64(data)
		}
	default:
		var data int64
		if err = binary.Read(r.Buf, binary.BigEndian, &data); err == nil {
			*value = int64(data)
		}
	}
	return
}

func (r *DataReader) Read(value interface{}) (err error) {
	T := reflect.TypeOf(value)
	if T.Kind() != reflect.Ptr {
		return errors.New("param must be pointer type")
	}

	V := reflect.ValueOf(value).Elem()
	return r._readItem(&V)
}

func (r *DataReader) _readItem(V *reflect.Value) (err error) {
	// data tag
	var tag DataTag
	if err = r.ReadTagData(&tag); err != nil {
		return
	}

	var child_sum uint64
	var struct_ver uint8

	switch V.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		if tag.DataType != TYPE_ARRAY {
			return errors.New("buf should be array type")
		}
		if (V.Kind() == reflect.Map) != (tag.VersionTag) {
			return errors.New("malformed array size tag")
		}
		if !V.CanSet() {
			return errors.New("param is not setable")
		}
		if err = r.ReadUintData(tag.SizeTag, &child_sum); err != nil {
			return
		}
	case reflect.Struct:
		if tag.DataType != TYPE_OBJECT {
			return errors.New("buf should be object type")
		}
		if tag.VersionTag {
			var tmpv uint64
			if err = r.ReadUintData(TAG_1, &tmpv); err != nil {
				return
			}
			struct_ver = uint8(tmpv)
		}
		if err = r.ReadUintData(tag.SizeTag, &child_sum); err != nil {
			return
		}
	case reflect.String:
		if err = r.ReadUintData(tag.SizeTag, &child_sum); err != nil {
			return
		}
	default:
		if !V.CanSet() {
			return errors.New("param is not setable")
		}
	}

	switch V.Kind() {
	case reflect.Map:
		V.Set(reflect.MakeMap(V.Type()))
		for i := 0; i < int(child_sum); i++ {
			mkey := reflect.New(V.Type().Key()).Elem()
			mValue := reflect.New(V.Type().Elem()).Elem()
			if err = r._readItem(&mkey); err != nil {
				return
			}
			if err = r._readItem(&mValue); err != nil {
				return
			}
			V.SetMapIndex(mkey, mValue)
		}

	case reflect.Slice:
		S := reflect.MakeSlice(V.Type(), 0, int(min(child_sum, 100)))
		for i := 0; i < int(child_sum); i++ {
			item := reflect.New(S.Type().Elem()).Elem()
			if err = r._readItem(&item); err != nil {
				return
			}
			S = reflect.Append(S, item)
		}
		V.Set(S)

	case reflect.Array:
		capacity := V.Len()
		if capacity != int(child_sum) {
			return errors.New("arrays' size should be equal")
		}
		for i := 0; i < int(child_sum); i++ {
			item := V.Index(i)
			r._readItem(&item)
		}

	case reflect.Struct:
		var curr_ver uint8
		real_childs := uint64(V.NumField())
		if tag.VersionTag {
			curr_ver = _getStructVersion(V)
			if (struct_ver >= curr_ver && child_sum < real_childs) || (struct_ver <= curr_ver && child_sum > real_childs) {
				return errors.New(fmt.Sprint(
					"mismatched version tag ",
					struct_ver, "/", child_sum,
					curr_ver, "/", real_childs))
			}
		}
		// read existed fields
		for i := 0; i < int(real_childs); i++ {
			item := V.Field(i)
			if err = r._readItem(&item); err != nil {
				return
			}
		}
		// ignore unsupport fields
		if tag.VersionTag {
			for i := 0; i < int(child_sum-real_childs); i++ {
				if err = r.ReadAndIgnore(); err != nil {
					return
				}
			}
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if tag.DataType != TYPE_INT {
			return errors.New("buf was not int type")
		}
		var tmpv int64
		if err = r.ReadIntData(tag.SizeTag, &tmpv); err != nil {
			return
		}
		V.SetInt(tmpv)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if tag.DataType != TYPE_UINT {
			return errors.New("buf was not uint type")
		}
		var tmpv uint64
		if err = r.ReadUintData(tag.SizeTag, &tmpv); err != nil {
			return
		}
		V.SetUint(tmpv)

	case reflect.Float32:
		if (tag.DataType != TYPE_REAL) || (tag.SizeTag != TAG_4) {
			return errors.New("buf was not float32 type")
		}
		var tmpv float32
		if err = binary.Read(r.Buf, binary.BigEndian, &tmpv); err != nil {
			return
		}
		V.SetFloat(float64(tmpv))

	case reflect.Float64:
		if (tag.DataType != TYPE_REAL) || (tag.SizeTag != TAG_8) {
			return errors.New("buf was not float64 type")
		}
		var tmpv float64
		if err = binary.Read(r.Buf, binary.BigEndian, &tmpv); err != nil {
			return
		}
		V.SetFloat(float64(tmpv))

	case reflect.Bool:
		if tag.DataType != TYPE_BOOL {
			return errors.New("buf was not bool type")
		}
		V.SetBool(tag.SizeTag == TAG_1)

	case reflect.String:
		if child_sum > 0 {
			tmps := make([]byte, child_sum)
			if err = binary.Read(r.Buf, binary.BigEndian, tmps); err != nil {
				return
			}
			V.SetString(string(tmps))
		} else {
			V.SetString("")
		}

	default:
		err = errors.New("unsupported read type " + V.Type().String())
	}
	return
}

func (r *DataReader) ReadAndIgnore() (err error) {
	var tag DataTag
	if err = r.ReadTagData(&tag); err != nil {
		return
	}

	var todo_bypes uint64
	switch tag.DataType {
	case TYPE_NONE, TYPE_BOOL:
	case TYPE_RAW, TYPE_STRING:
		// string & raw read len, and read data.
		if err = r.ReadUintData(tag.SizeTag, &todo_bypes); err != nil {
			return
		}
		break
	case TYPE_INT, TYPE_UINT, TYPE_REAL:
		switch tag.SizeTag {
		case TAG_0:
		case TAG_1:
			todo_bypes = 1
		case TAG_2:
			todo_bypes = 2
		case TAG_4:
			todo_bypes = 4
		default:
			todo_bypes = 8
		}

	case TYPE_ARRAY:
		// array & object -> read child_sum, ver, and contents.
		var child_sum uint64
		if err = r.ReadUintData(tag.SizeTag, &child_sum); err != nil {
			return
		}

		if tag.VersionTag {
			for i := 0; i < int(child_sum); i++ {
				if err = r.ReadAndIgnore(); err != nil {
					return
				}
			}
		} else {
			for i := 0; i < int(child_sum); i++ {
				if err = r.ReadAndIgnore(); err != nil {
					return
				}
				if err = r.ReadAndIgnore(); err != nil {
					return
				}
			}
		}

	case TYPE_OBJECT:
		var child_sum uint64
		if err = r.ReadUintData(tag.SizeTag, &child_sum); err != nil {
			return
		}

		if tag.VersionTag {
			var tmpv uint8
			if err = binary.Read(r.Buf, binary.BigEndian, &tmpv); err != nil {
				return
			}
		}

		for i := 0; i < int(child_sum); i++ {
			if err = r.ReadAndIgnore(); err != nil {
				return
			}
		}

	default:
		return errors.New("unknown data type")
	}

	return binary.Read(r.Buf, binary.BigEndian, make([]byte, todo_bypes, todo_bypes))
}
