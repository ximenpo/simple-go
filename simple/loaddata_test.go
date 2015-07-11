package simple

import (
	"testing"
)

func Test_LoadData_FromFile(t *testing.T) {
	if data, err := LoadData_FromFile("./loaddata_test.txt"); err != nil {
		t.Error(err)
	} else {
		t.Log(string(data))
	}
}
