package utils

// serialise/deserialise a bunch of variables
type Shifter struct {
	values         map[int]*Value
	buf            []byte
	consumed_bytes int
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
	if len(sh.buf) < sh.consumed_bytes {
		return 0
	}
	res := sh.buf[sh.consumed_bytes]
	sh.consumed_bytes++
	return res
}
