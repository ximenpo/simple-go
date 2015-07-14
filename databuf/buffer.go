package databuf

import (
	"errors"
	"fmt"
	"strings"
)

var (
	Buffer_DefaultCapacity uint = 4 * 1024
)

type Buffer struct {
	data    []byte
	capcity uint
	size    uint
	offset  uint
	failure bool
}

func NewBuffer(data []byte) (buf *Buffer) {
	buf = &Buffer{}
	if data != nil {
		buf.Assign(data)
	}
	return
}

func (self *Buffer) Assign(data []byte) {
	self.capcity = uint(len(data))
	self.size = self.capcity
	self.data = data
	self.offset = 0
	self.failure = false
}

func (self *Buffer) Clear() {
	self.failure = false
	self.size = 0
	self.offset = 0
}

func (self *Buffer) Rewind() {
	self.failure = false
	self.offset = 0
}

func (self *Buffer) Data() []byte {
	return self.data[0:self.size]
}

func (self *Buffer) Size() uint {
	return self.size
}

func (self *Buffer) Capacity() uint {
	return self.capcity
}

func (self *Buffer) Pos() uint {
	return self.offset
}

func (self *Buffer) DataFromCurrPos() []byte {
	return self.data[self.offset:self.size]
}

func (self *Buffer) Good() bool {
	return !self.failure
}

func (self *Buffer) Failure() bool {
	return self.failure
}

func (self *Buffer) SetFailure() {
	self.failure = true
}

func (self *Buffer) Write(pData []byte) (n int, err error) {
	if self.Failure() {
		return 0, errors.New("buffer was failure")
	}

	if pData == nil {
		return 0, nil
	}

	if self.capcity == 0 {
		self.data = make([]byte, Buffer_DefaultCapacity)
		self.capcity = Buffer_DefaultCapacity
	}

	nLen := uint(len(pData))
	if self.offset+nLen > self.capcity {
		self.failure = true
		return 0, errors.New("now enough space")
	}

	copy(self.data[self.offset:], pData[0:nLen])
	self.offset += nLen
	if self.offset > self.size {
		self.size = self.offset
	}
	return int(nLen), nil
}

func (self *Buffer) Read(pData []byte) (n int, err error) {
	if self.Failure() {
		return 0, errors.New("buffer was failure")
	}

	if pData == nil {
		return 0, nil
	}

	nLen := uint(len(pData))
	if self.size-self.offset < nLen {
		self.failure = true
		return 0, errors.New("now enough data")
	}

	copy(pData, self.data[self.offset:self.offset+nLen])
	self.offset += nLen
	return int(nLen), nil
}

func (self *Buffer) Dump() string {
	line := ""
	for i := uint(0); i < self.Size(); i++ {
		line += fmt.Sprintf("%02X", self.data[i])
		line += " "
		if i%16 == 15 {
			line += "\n"
		}
	}
	return strings.ToUpper(line)
}
