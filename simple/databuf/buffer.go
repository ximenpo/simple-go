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
	failure bool
	data    []byte
	capcity uint
	size    uint
	offset  uint
}

func (b *Buffer) Assign(data []byte) {
	b.capcity = uint(len(data))
	b.size = b.capcity
	b.data = data
	b.offset = 0
	b.failure = false
}

func (b *Buffer) Clear() {
	b.failure = false
	b.size = 0
	b.offset = 0
}

func (b *Buffer) Rewind() {
	b.failure = false
	b.offset = 0
}

func (b *Buffer) Data() []byte {
	return b.data[0:b.size]
}

func (b *Buffer) Size() uint {
	return b.size
}

func (b *Buffer) Capacity() uint {
	return b.capcity
}

func (b *Buffer) Pos() uint {
	return b.offset
}

func (b *Buffer) DataFromCurrPos() []byte {
	return b.data[b.offset:b.size]
}

func (b *Buffer) Good() bool {
	return !b.failure
}

func (b *Buffer) Failure() bool {
	return b.failure
}

func (b *Buffer) SetFailure() {
	b.failure = true
}

func (b *Buffer) Write(pData []byte) (n int, err error) {
	if b.Failure() {
		return 0, errors.New("buffer was failure")
	}

	if pData == nil {
		return 0, nil
	}

	if b.capcity == 0 {
		b.data = make([]byte, Buffer_DefaultCapacity)
		b.capcity = Buffer_DefaultCapacity
	}

	nLen := uint(len(pData))
	if b.offset+nLen > b.capcity {
		b.failure = true
		return 0, errors.New("now enough space")
	}

	copy(b.data[b.offset:], pData[0:nLen])
	b.offset += nLen
	if b.offset > b.size {
		b.size = b.offset
	}
	return int(nLen), nil
}

func (b *Buffer) Read(pData []byte) (n int, err error) {
	if b.Failure() {
		return 0, errors.New("buffer was failure")
	}

	if pData == nil {
		return 0, nil
	}

	nLen := uint(len(pData))
	if b.size-b.offset < nLen {
		b.failure = true
		return 0, errors.New("now enough data")
	}

	copy(pData, b.data[b.offset:b.offset+nLen])
	b.offset += nLen
	return int(nLen), nil
}

func (b *Buffer) Dump() string {
	line := ""
	for i := uint(0); i < b.Size(); i++ {
		line += fmt.Sprintf("%02X", b.data[i])
		line += " "
		if i%16 == 15 {
			line += "\n"
		}
	}
	return strings.ToUpper(line)
}
