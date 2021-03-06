package databuf

import (
	"math"
	"reflect"
)

// DATA_TYPE
const (
	TYPE_NONE   = iota // nothing, dummy data
	TYPE_RAW           // byte array
	TYPE_BOOL          // boolean
	TYPE_INT           // int8, int16, int32
	TYPE_UINT          // uint8, uint16, uint32
	TYPE_REAL          // float32, float64
	TYPE_STRING        // UTF-8 strings
	TYPE_ARRAY         // array
	TYPE_OBJECT        // key/value pairs
	TYPE_SUM
)

// SIZE_TAG
const (
	TAG_0 = iota // 0byte
	TAG_1        // 1byte
	TAG_2        // 2byte
	TAG_4        // 4byte
	TAG_8        // 8byte
	TAG_SUM
)

type DataTag struct {
	DataType   uint
	SizeTag    uint
	VersionTag bool
}

func (self DataTag) Pack() (ret uint8) {
	ret = uint8((self.SizeTag << 4) | (self.DataType << 0))
	if self.VersionTag {
		ret |= 0x80
	}
	return
}

func (self *DataTag) UnPack(value uint8) {
	self.DataType = uint(value & 0x0F)
	self.SizeTag = uint((value & 0x70) >> 4)
	self.VersionTag = (value & 0x80) != 0
}

func DataSizeTag(value interface{}) (ret uint) {
	ret = TAG_0
	switch v := value.(type) {
	case bool:
		if v {
			ret = TAG_1
		} else {
			ret = TAG_0
		}
	case int, int8, int16, int32, int64:
		vi := reflect.ValueOf(v).Int()
		switch {
		case vi == 0:
			ret = TAG_0
		case vi <= math.MaxInt8:
			ret = TAG_1
		case vi <= math.MaxInt16:
			ret = TAG_2
		case vi <= math.MaxInt32:
			ret = TAG_4
		default:
			ret = TAG_8
		}
	case uint, uint8, uint16, uint32, uint64:
		vi := reflect.ValueOf(v).Uint()
		switch {
		case vi == 0:
			ret = TAG_0
		case vi <= math.MaxUint8:
			ret = TAG_1
		case vi <= math.MaxUint16:
			ret = TAG_2
		case vi <= math.MaxUint32:
			ret = TAG_4
		default:
			ret = TAG_8
		}
	case float32:
		return TAG_4
	case float64:
		return TAG_8
	default:
		panic("unsupport tag size type:" + reflect.TypeOf(value).String())
	}
	return
}
