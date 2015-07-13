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
	Buf io.Writer
}

func NewDataWriter(buf io.Writer) *DataWriter {
	return &DataWriter{buf}
}

func (self *DataWriter) WriteTagData(tag DataTag) (err error) {
	return binary.Write(self.Buf, binary.BigEndian, tag.Pack())
}

func (self *DataWriter) WriteUintData(size_tag uint, value uint64) (err error) {
	switch size_tag {
	case TAG_0:
		return nil
	case TAG_1:
		return binary.Write(self.Buf, binary.BigEndian, uint8(value))
	case TAG_2:
		return binary.Write(self.Buf, binary.BigEndian, uint16(value))
	case TAG_4:
		return binary.Write(self.Buf, binary.BigEndian, uint32(value))
	default:
		return binary.Write(self.Buf, binary.BigEndian, uint64(value))
	}
}

func (self *DataWriter) WriteIntData(size_tag uint, value int64) (err error) {
	switch size_tag {
	case TAG_0:
		return nil
	case TAG_1:
		return binary.Write(self.Buf, binary.BigEndian, int8(value))
	case TAG_2:
		return binary.Write(self.Buf, binary.BigEndian, int16(value))
	case TAG_4:
		return binary.Write(self.Buf, binary.BigEndian, int32(value))
	default:
		return binary.Write(self.Buf, binary.BigEndian, int64(value))
	}
}

func (self *DataWriter) Write(value interface{}) (err error) {
	V := reflect.ValueOf(value)
	if V.Kind() == reflect.Ptr {
		return self.Write(V.Elem().Interface())
	}

	var tag DataTag
	switch V.Kind() {
	case reflect.Map:
		keys := V.MapKeys()
		length := len(keys)
		tag := DataTag{TYPE_ARRAY, DataSizeTag(length), true}
		if err = self.WriteTagData(tag); err != nil {
			return
		}
		if err = self.WriteUintData(tag.SizeTag, uint64(length)); err != nil {
			return
		}
		for i := 0; i < length; i++ {
			v := V.MapIndex(keys[i])
			if err = self.Write(keys[i].Interface()); err != nil {
				return
			}
			if err = self.Write(v.Interface()); err != nil {
				return
			}
		}

	case reflect.Slice, reflect.Array:
		length := V.Len()
		tag := DataTag{TYPE_ARRAY, DataSizeTag(length), false}
		if err = self.WriteTagData(tag); err != nil {
			return
		}
		if err = self.WriteUintData(tag.SizeTag, uint64(length)); err != nil {
			return
		}
		for i := 0; i < length; i++ {
			v := V.Index(i)
			if err = self.Write(v.Interface()); err != nil {
				return
			}
		}

	case reflect.Struct:
		fields := V.NumField()
		curr_ver := _getStructVersion(&V)
		tag := DataTag{TYPE_OBJECT, DataSizeTag(fields), curr_ver != 0}
		if err = self.WriteTagData(tag); err != nil {
			return
		}
		if tag.VersionTag {
			if err = self.WriteUintData(TAG_1, uint64(curr_ver)); err != nil {
				return
			}
		}
		if err = self.WriteUintData(tag.SizeTag, uint64(fields)); err != nil {
			return
		}
		for i := 0; i < fields; i++ {
			f := V.Field(i)
			if err = self.Write(f.Interface()); err != nil {
				return
			}
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		tag := DataTag{TYPE_INT, DataSizeTag(value), false}
		if err = self.WriteTagData(tag); err != nil {
			return
		}
		return self.WriteIntData(tag.SizeTag, V.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		tag := DataTag{TYPE_UINT, DataSizeTag(value), false}
		if err = self.WriteTagData(tag); err != nil {
			return
		}
		return self.WriteUintData(tag.SizeTag, V.Uint())

	case reflect.Float32, reflect.Float64:
		if V.Kind() == reflect.Float32 {
			tag = DataTag{TYPE_REAL, TAG_4, false}
		} else {
			tag = DataTag{TYPE_REAL, TAG_8, false}
		}
		if err = self.WriteTagData(tag); err != nil {
			return
		}
		return binary.Write(self.Buf, binary.BigEndian, value)

	case reflect.Bool:
		if V.Bool() {
			tag = DataTag{TYPE_BOOL, TAG_1, false}
		} else {
			tag = DataTag{TYPE_BOOL, TAG_0, false}
		}
		if err = self.WriteTagData(tag); err != nil {
			return
		}

	case reflect.String:
		cBuff := []byte(V.String())
		tag = DataTag{TYPE_STRING, DataSizeTag(len(cBuff)), false}
		if err = self.WriteTagData(tag); err != nil {
			return
		}
		if err = self.WriteUintData(tag.SizeTag, uint64(len(cBuff))); err != nil {
			return
		}
		if _, err = self.Buf.Write([]byte(cBuff)); err != nil {
			return
		}

	default:
		return errors.New("unsupported type " + reflect.TypeOf(value).String())
	}

	return
}
