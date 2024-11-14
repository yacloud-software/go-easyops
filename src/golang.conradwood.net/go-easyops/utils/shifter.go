package utils

import (
	"fmt"
)

// serialise/deserialise a bunch of variables
type Shifter struct {
	values         map[int]*Value
	buf            []byte
	consumed_bytes int
	err            error
}
type Value struct {
	Integer int
}

/*
A "Shifter" shifts bytes from a larger number into an array and unshifts it again
*/
func NewShifter(buf []byte) *Shifter {
	res := &Shifter{
		buf:            buf,
		consumed_bytes: 0,
	}
	return res
}

// next byte is length, followed by uint8s
func (sh *Shifter) Array8() []byte {
	l := int(sh.Unshift_uint8())
	res := make([]byte, l)
	for i := 0; i < l; i++ {
		res[i] = sh.Unshift_uint8()
	}
	return res
}

func (sh *Shifter) Unshift_uint16() uint32 {
	res := uint32(0)
	for i := 0; i < 2; i++ {
		b := uint32(sh.Unshift_uint8())
		b = b << (8 * i)
		res = res | b
	}
	return res
}

// LSB first, MSB last. e.g. 0xAABBCCDD will be shifted into a byte array like so: []byte{0xDD,0xCC,0xBB,0xAA}
func (sh *Shifter) Unshift_uint32() uint32 {
	res := uint32(0)
	for i := 0; i < 4; i++ {
		b := uint32(sh.Unshift_uint8())
		b = b << (8 * i)
		res = res | b
	}
	return res
}
func (sh *Shifter) Unshift_uint64() uint64 {
	res := uint64(0)
	for i := 0; i < 8; i++ {
		b := uint64(sh.Unshift_uint8())
		b = b << (8 * i)
		res = res | b
	}
	return res
}
func (sh *Shifter) Unshift_uint8() uint8 {
	if len(sh.buf) <= sh.consumed_bytes {
		sh.consumed_bytes++
		sh.err = fmt.Errorf("Read beyond length (length=%d, read=%d)", len(sh.buf), sh.consumed_bytes)
		return 0
	}
	res := sh.buf[sh.consumed_bytes]
	sh.consumed_bytes++
	return res
}

// modify the underlying byte array to set a uint32 at a particular address
func (sh *Shifter) SetUint32(pos int, b uint32) {
	if pos+4 >= len(sh.buf) {
		sh.err = fmt.Errorf("Write beyond length (length=%d, write=%d)", len(sh.buf), pos)
		return
	}
	sh.buf[pos] = byte(0xFF & b)
	sh.buf[pos+1] = byte(0xFF & (b >> 8))
	sh.buf[pos+2] = byte(0xFF & (b >> 16))
	sh.buf[pos+3] = byte(0xFF & (b >> 24))
}
func (sh *Shifter) Bytes() []byte {
	return sh.buf
}
func (sh *Shifter) Error() error {
	return sh.err
}
