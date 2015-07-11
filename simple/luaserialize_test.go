package simple

import (
	"testing"
)

func Test_SerializeToLuaCode(t *testing.T) {
	var b bool
	if ret, err := LuaSerialize(&b); (err != nil) || (ret != "false") {
		t.Error(ret, err)
	}
}
