package simple

import (
	"testing"
)

type TestIntfA interface {
	GetName() string
}

type TestIntfB interface {
	GetAge() int
}

type TestS struct {
}

func (self *TestS) GetName() string {
	return "person.name"
}
func (self *TestS) GetAge() int {
	return 35
}

type TestC struct {
	TestIntfA
	TestIntfB
}

func Test_InfectFieldsByType(t *testing.T) {
	obj := TestC{}

	s := TestS{}
	var a TestIntfA = &s
	sum, err := InfectFields(&obj, &a)
	if err != nil {
		t.Error(err)
	}
	if sum != 1 {
		t.Error("wrong infect fields sum")
	}
}
