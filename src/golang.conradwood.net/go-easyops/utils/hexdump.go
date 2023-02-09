package utils

import (
	"bytes"
	"fmt"
	"unicode"
)

// return a mini hexdump as single line string
func HexStr(buf []byte) string {
	maxlen := 24
	if len(buf) > maxlen {
		buf = buf[:maxlen]
	}
	s := ""
	deli := ""
	for _, b := range buf {
		s = s + deli + fmt.Sprintf("%02X", b)
		deli = " "
	}
	s = s + " "
	x := string(buf)
	for _, r := range x {
		if unicode.IsPrint(r) {
			s = s + string(r)
		} else {
			s = s + "."
		}
	}
	return s
}

// return the buffer as a hexdump with ascii. each line prefixed by 'prefix'
func Hexdump(prefix string, buf []byte) string {
	return HexdumpWithLen(60, prefix, buf)
}
func HexdumpWithLen(ln int, prefix string, buf []byte) string {
	var res bytes.Buffer
	linelen := 0
	needprefix := true
	ascii := ""
	for _, b := range buf {
		if needprefix {
			res.WriteString(prefix)
		}
		needprefix = false
		s := fmt.Sprintf("%02X ", b)
		ascii = ascii + hextochar(b)
		res.WriteString(s)
		linelen = linelen + len(s)
		if linelen >= ln {
			linelen = 0
			res.WriteString(ascii)
			ascii = ""
			res.WriteString("\n")
			needprefix = true
		}
	}
	if len(ascii) != 0 {
		res.WriteString(ascii)
		res.WriteString("\n")
	}
	return res.String()
}
func hextochar(b byte) string {
	r := rune(b)
	if !unicode.IsPrint(r) {
		return "."
	}
	return string(r)

}
