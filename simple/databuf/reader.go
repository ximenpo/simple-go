package databuf

import (
	"encoding/binary"
	"errors"
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
	buf io.Reader
}

func (r *DataReader) ReadTagData(tag *DataTag) (err error) {
	var data uint8
	if err = binary.Read(r.buf, binary.BigEndian, &data); err == nil {
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
		if err = binary.Read(r.buf, binary.BigEndian, &data); err == nil {
			*value = uint64(data)
		}
	case TAG_2:
		var data uint16
		if err = binary.Read(r.buf, binary.BigEndian, &data); err == nil {
			*value = uint64(data)
		}
	case TAG_4:
		var data uint32
		if err = binary.Read(r.buf, binary.BigEndian, &data); err == nil {
			*value = uint64(data)
		}
	default:
		var data uint64
		if err = binary.Read(r.buf, binary.BigEndian, &data); err == nil {
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
		if err = binary.Read(r.buf, binary.BigEndian, &data); err == nil {
			*value = int64(data)
		}
	case TAG_2:
		var data int16
		if err = binary.Read(r.buf, binary.BigEndian, &data); err == nil {
			*value = int64(data)
		}
	case TAG_4:
		var data int32
		if err = binary.Read(r.buf, binary.BigEndian, &data); err == nil {
			*value = int64(data)
		}
	default:
		var data int64
		if err = binary.Read(r.buf, binary.BigEndian, &data); err == nil {
			*value = int64(data)
		}
	}
	return
}

func (r *DataReader) Read(value interface{}) (err error) {
	{
		T := reflect.TypeOf(value)
		if T.Kind() != reflect.Ptr {
			return errors.New("param must be pointer type")
		}
	}

	// data tag
	var tag DataTag
	if err = r.ReadTagData(&tag); err != nil {
		return
	}

	var size uint64

	V := reflect.ValueOf(value).Elem()
	switch V.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		if tag.DataType != TYPE_ARRAY {
			return errors.New("buf should be array type")
		}
		if (V.Kind() == reflect.Map) == (tag.VersionTag) {
			return errors.New("malformed array size tag")
		}
		if !V.CanSet() {
			return errors.New("param is not setable")
		}
		if err = r.ReadUintData(tag.SizeTag, &size); err != nil {
			return
		}
	case reflect.Struct:
		if tag.DataType != TYPE_OBJECT {
			return errors.New("buf should be object type")
		}
		if err = r.ReadUintData(tag.SizeTag, &size); err != nil {
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
		for i := 0; i < int(size); i++ {
			mkey := reflect.New(V.Type().Key()).Elem()
			mValue := reflect.New(V.Type().Elem()).Elem()
			if err = r.Read(&mkey); err != nil {
				return
			}
			if err = r.Read(&mValue); err != nil {
				return
			}
			V.SetMapIndex(mkey, mValue)
		}

	case reflect.Slice:
		S := reflect.MakeSlice(V.Type(), 0, int(min(size, 100)))
		for i := 0; i < int(size); i++ {
			item := reflect.New(S.Type().Elem()).Elem()
			if err = r.Read(&item); err != nil {
				return
			}
			S = reflect.Append(S, item)
		}
		V.Set(S)

	case reflect.Array:
		capacity := V.Len()
		if capacity != int(size) {
			return errors.New("arrays' size should be equal")
		}
		for i := 0; i < int(size); i++ {
			item := V.Index(i)
			r.Read(&item)
		}

	case reflect.Struct:
		var ver, curr_ver uint8
		if tag.VersionTag {
			var tmpv uint64
			if err = r.ReadUintData(TAG_1, &tmpv); err != nil {
				return
			}
			ver = uint8(tmpv)
			curr_ver = _getStructVersion(&V)
		}
		realSize := uint64(V.NumField())
		if (tag.VersionTag && ver >= curr_ver && size < realSize) ||
			(tag.VersionTag && ver <= curr_ver && size > realSize) {
			return errors.New("mismatched version tag")
		}
		// read existed fields
		for i := 0; i < int(realSize); i++ {
			item := V.Field(i)
			if err = r.Read(&item); err != nil {
				return
			}
		}
		// ignore unsupport fields
		if tag.VersionTag {
			for i := 0; i < int(size-realSize); i++ {
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
		if err = binary.Read(r.buf, binary.BigEndian, &tmpv); err != nil {
			return
		}
		V.SetFloat(float64(tmpv))

	case reflect.Float64:
		if (tag.DataType != TYPE_REAL) || (tag.SizeTag != TAG_8) {
			return errors.New("buf was not float64 type")
		}
		var tmpv float64
		if err = binary.Read(r.buf, binary.BigEndian, &tmpv); err != nil {
			return
		}
		V.SetFloat(float64(tmpv))

	case reflect.Bool:
		if tag.DataType != TYPE_BOOL {
			return errors.New("buf was not bool type")
		}
		V.SetBool(tag.SizeTag == TAG_1)

	case reflect.String:
		if size > 0 {
			tmps := make([]byte, size)
			if err = binary.Read(r.buf, binary.BigEndian, tmps); err != nil {
				return
			}
			V.SetString(string(tmps))
		} else {
			V.SetString("")
		}

	default:
		err = errors.New("unsupported read type " + reflect.TypeOf(value).String())
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
		// array & object -> read size, ver, and contents.
		var size uint64
		if err = r.ReadUintData(tag.SizeTag, &size); err != nil {
			return
		}

		if tag.VersionTag {
			for i := 0; i < int(size); i++ {
				if err = r.ReadAndIgnore(); err != nil {
					return
				}
			}
		} else {
			for i := 0; i < int(size); i++ {
				if err = r.ReadAndIgnore(); err != nil {
					return
				}
				if err = r.ReadAndIgnore(); err != nil {
					return
				}
			}
		}

	case TYPE_OBJECT:
		var size uint64
		if err = r.ReadUintData(tag.SizeTag, &size); err != nil {
			return
		}

		if tag.VersionTag {
			var tmpv uint8
			if err = binary.Read(r.buf, binary.BigEndian, &tmpv); err != nil {
				return
			}
		}

		for i := 0; i < int(size); i++ {
			if err = r.ReadAndIgnore(); err != nil {
				return
			}
		}

	default:
		return errors.New("unknown data type")
	}

	return binary.Read(r.buf, binary.BigEndian, make([]byte, todo_bypes, todo_bypes))
}
