package utils

import (
	"bytes"
	"fmt"
	"unicode"
)

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
