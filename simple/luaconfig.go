package simple

import (
	"errors"
	"fmt"
	"github.com/hoisie/mustache"
	"github.com/ximenpo/simple-luago/lua"
	"unsafe"
)

func LuaFetchConfig(addr string, config interface{},
	config_field_must_exist bool, pre_defines ...interface{}) (err error) {
	// fetch data
	real_addr := mustache.Render(addr, pre_defines...)
	data, err := LoadData(real_addr)
	if err != nil {
		return
	}

	vm := lua.LuaVM{}
	vm.Start()
	defer vm.Stop()

	// set predeined value
	for i := range pre_defines {
		define := pre_defines[i]
		if !vm.SetObject("_G", define, true) {
			return errors.New(fmt.Sprint("set lua define variable failed for ", i))
		}
	}

	if err = vm.RunBuffer(unsafe.Pointer(&data[0]), uint(len(data))); err != nil {
		return
	}

	if !vm.GetObject("_G", config, !config_field_must_exist) {
		return errors.New("get lua config variable error ")
	}

	return nil
}
