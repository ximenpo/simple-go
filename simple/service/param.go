package service

import (
	"errors"
	"flag"
)

type Params struct {
	Type   string // service type
	Group  string // service group, platform ... etc.
	ID     string // service id
	Config string // service config file path/url
}

func (self Params) String() string {
	return self.Type + "/" + self.Group + "/" + self.ID
}

func (self Params) Valid() bool {
	return (self.Type != "") || (self.ID != "")
}

func (self *Params) LoadFromFlag() error {
	if !flag.Parsed() {
		return errors.New("flag was not parsed")
	}

	*self = _cmdlineParams
	return nil
}

var _cmdlineParams Params

func init() {
	flag.StringVar(&_cmdlineParams.Type, "type", "", "service type")
	flag.StringVar(&_cmdlineParams.Group, "group", "", "service type")
	flag.StringVar(&_cmdlineParams.ID, "id", "", "service type")
	flag.StringVar(&_cmdlineParams.Config, "config", "", "service type")
}
