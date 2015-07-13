package databuf

import (
	"testing"
)

func Test_SizeTag(t *testing.T) {
	if TAG_0 != DataSizeTag(0) {
		t.Error("TAG_0 != DataSizeTag(0)")
	}
	if TAG_1 != DataSizeTag(13) {
		t.Errorf("TAG_1 != %d", DataSizeTag(13))
	}
}
