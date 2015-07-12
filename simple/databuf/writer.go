package databuf

import (
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"strconv"
)

func _getStructVersion(V *reflect.Value) uint8 {
	var ver uint64
	for i := 0; i < V.NumField(); i++ {
		fieldVersion := V.Type().Field(i).Tag.Get("version")
		if 0 == len(fieldVersion) {
			continue
		}

		if t, err := strconv.ParseUint(fieldVersion, 10, 64); (err == nil) && (t > ver) {
			ver = t
		}
	}

	return uint8(ver)
}

type DataWriter struct {
	buf io.Writer
}

func NewWriter(buf io.Writer) *DataWriter {
	return &DataWriter{buf}
}

func (w *DataWriter) WriteTagData(tag DataTag) (err error) {
	return binary.Write(w.buf, binary.BigEndian, tag.Pack())
}

func (w *DataWriter) WriteUintData(size_tag uint, value uint64) (err error) {
	switch size_tag {
	case TAG_0:
		return nil
	case TAG_1:
		return binary.Write(w.buf, binary.BigEndian, uint8(value))
	case TAG_2:
		return binary.Write(w.buf, binary.BigEndian, uint16(value))
	case TAG_4:
		return binary.Write(w.buf, binary.BigEndian, uint32(value))
	default:
		return binary.Write(w.buf, binary.BigEndian, uint64(value))
	}
}

func (w *DataWriter) WriteIntData(size_tag uint, value int64) (err error) {
	switch size_tag {
	case TAG_0:
		return nil
	case TAG_1:
		return binary.Write(w.buf, binary.BigEndian, int8(value))
	case TAG_2:
		return binary.Write(w.buf, binary.BigEndian, int16(value))
	case TAG_4:
		return binary.Write(w.buf, binary.BigEndian, int32(value))
	default:
		return binary.Write(w.buf, binary.BigEndian, int64(value))
	}
}

func (w *DataWriter) Write(value interface{}) (err error) {
	V := reflect.ValueOf(value)
	if V.Kind() == reflect.Ptr {
		return w.Write(V.Elem().Interface())
	}

	var tag DataTag
	switch V.Kind() {
	case reflect.Map:
		keys := V.MapKeys()
		length := len(keys)
		tag := DataTag{TYPE_ARRAY, SizeTag(length), false}
		if err = w.WriteTagData(tag); err != nil {
			return
		}
		if err = w.WriteUintData(tag.SizeTag, uint64(length)); err != nil {
			return
		}
		for i := 0; i < length; i++ {
			v := V.MapIndex(keys[i])
			if err = w.Write(keys[i].Interface()); err != nil {
				return
			}
			if err = w.Write(v.Interface()); err != nil {
				return
			}
		}

	case reflect.Slice, reflect.Array:
		length := V.Len()
		tag := DataTag{TYPE_ARRAY, SizeTag(length), false}
		if err = w.WriteTagData(tag); err != nil {
			return
		}
		if err = w.WriteUintData(tag.SizeTag, uint64(length)); err != nil {
			return
		}
		for i := 0; i < length; i++ {
			v := V.Index(i)
			if err = w.Write(v.Interface()); err != nil {
				return
			}
		}

	case reflect.Struct:
		length := V.NumField()
		stVersion := _getStructVersion(&V)
		tag := DataTag{TYPE_OBJECT, SizeTag(length), stVersion != 0}
		if err = w.WriteTagData(tag); err != nil {
			return
		}
		if err = w.WriteUintData(tag.SizeTag, uint64(length)); err != nil {
			return
		}
		if tag.VersionTag {
			if err = w.WriteUintData(TAG_1, uint64(stVersion)); err != nil {
				return
			}
		}
		for i := 0; i < length; i++ {
			v := V.Field(i)
			if err = w.Write(v.Interface()); err != nil {
				return
			}
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		tag := DataTag{TYPE_INT, SizeTag(value), false}
		if err = w.WriteTagData(tag); err != nil {
			return
		}
		return w.WriteIntData(tag.SizeTag, V.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		tag := DataTag{TYPE_UINT, SizeTag(value), false}
		if err = w.WriteTagData(tag); err != nil {
			return
		}
		return w.WriteUintData(tag.SizeTag, V.Uint())

	case reflect.Float32, reflect.Float64:
		if V.Kind() == reflect.Float32 {
			tag = DataTag{TYPE_REAL, TAG_4, false}
		} else {
			tag = DataTag{TYPE_REAL, TAG_8, false}
		}
		if err = w.WriteTagData(tag); err != nil {
			return
		}
		return binary.Write(w.buf, binary.BigEndian, value)

	case reflect.Bool:
		if V.Bool() {
			tag = DataTag{TYPE_BOOL, TAG_1, false}
		} else {
			tag = DataTag{TYPE_BOOL, TAG_0, false}
		}
		if err = w.WriteTagData(tag); err != nil {
			return
		}

	case reflect.String:
		cBuff := []byte(V.String())
		tag = DataTag{TYPE_STRING, SizeTag(len(cBuff)), false}
		if err = w.WriteTagData(tag); err != nil {
			return
		}
		if err = w.WriteUintData(tag.SizeTag, uint64(len(cBuff))); err != nil {
			return
		}
		if _, err = w.buf.Write([]byte(cBuff)); err != nil {
			return
		}

	default:
		return errors.New("unsupported type " + reflect.TypeOf(value).String())
	}

	return
}
