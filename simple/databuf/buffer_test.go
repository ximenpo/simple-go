package databuf

import (
	"encoding/json"
	"reflect"
	"testing"
)

type TestStruct struct {
	V1  int8
	V2  uint8
	V3  int16
	V4  uint16
	V5  int32
	V6  uint32
	V7  int64
	V8  uint64
	V9  float32
	V10 float64
	V11 string
	V12 [2][2]int
	V13 map[int]string `version:"1"`
}

func Test_BufferSerialize(t *testing.T) {
	var testS, testT TestStruct
	testS.V1 = 1
	testS.V2 = 254
	testS.V3 = -2
	testS.V4 = 258
	testS.V5 = 555
	testS.V6 = 65538
	testS.V7 = 332
	testS.V8 = 5555
	testS.V9 = 1.414
	testS.V10 = 3.1415926
	testS.V11 = "测试程序test壹IIぁアЙá"
	testS.V12 = [2][2]int{{1, 2}, {55, -7}}
	testS.V13 = make(map[int]string)
	testS.V13[1] = "Test1"
	testS.V13[2] = "Test1"
	testS.V13[3] = "Test1"
	testS.V13[4] = "Test1"

	var buf Buffer
	{
		writer := NewDataWriter(&buf)
		if err := writer.Write(&testS); err != nil {
			t.Error(err)
		}
	}

	buf.Rewind()
	{
		reader := NewDataReader(&buf)
		if err := reader.Read(&testT); err != nil {
			t.Error(err)
		}
	}

	if !reflect.DeepEqual(testT, testS) {
		t.Logf("buf status: cap %d, size %d, pos %d", buf.Capacity(), buf.Size(), buf.Pos())
		if b, e := json.Marshal(testT); e == nil {
			t.Log(string(b))
		}
		t.Error("Read: \n", buf.Dump())
	}
}
