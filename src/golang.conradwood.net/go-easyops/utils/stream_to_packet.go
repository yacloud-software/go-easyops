package utils

import (
	"bytes"
	"fmt"
)

type StreamToPacket struct {
	buf              *bytes.Buffer
	start            byte
	esc              byte
	stop             byte
	last_was_escaped bool
	inpacket         bool
	debug            bool
}

func NewStreamToPacket(start, esc, stop byte) (*StreamToPacket, error) {
	res := &StreamToPacket{
		start: start,
		esc:   esc,
		debug: *debug_packetizer,
		stop:  stop,
		buf:   &bytes.Buffer{},
	}
	res.debugf("New stream packet, start=%c, esc=%c, stop=%c\n", start, esc, stop)
	return res, nil
}
func bdis(b byte) string {
	return fmt.Sprintf("0x%02X (%c)", b, b)
}

// returns true if a complete packet is in the buf
func (stp *StreamToPacket) AddByte(b byte) bool {
	stp.debugf("adding byte %s (inpacket=%v,last-escaped=%v)\n", bdis(b), stp.inpacket, stp.last_was_escaped)
	if !stp.inpacket {
		// scan for beginning of packet
		if stp.last_was_escaped {
			stp.last_was_escaped = false
			return false
		}
		if b == stp.esc {
			stp.last_was_escaped = true
			return false
		}
		if b != stp.start {
			return false
		}
		stp.inpacket = true
		return false
	}
	// we are in a packet, scan for a new beginning, then we need to restart/resynchronise
	if !stp.last_was_escaped && b == stp.start {
		stp.buf.Reset()
		return false
	}
	if stp.last_was_escaped {
		stp.buf.Write([]byte{b})
		stp.last_was_escaped = false
		return false
	}
	if b == stp.esc {
		stp.last_was_escaped = true
		return false
	}
	// scan for end of packet
	if b == stp.stop {
		return true
	}
	stp.buf.Write([]byte{b})
	return false
}
func (stp *StreamToPacket) ReadPacket() []byte {
	bt := stp.buf.Bytes()
	res := make([]byte, len(bt))
	for i, _ := range bt {
		res[i] = bt[i]
	}
	stp.buf.Reset()
	stp.last_was_escaped = false
	stp.inpacket = false
	stp.debugf("packet read\n")
	return res
}

func (stp *StreamToPacket) debugf(format string, args ...interface{}) {
	if !stp.debug {
		return
	}
	s := fmt.Sprintf(format, args...)
	fmt.Printf("[streamtopacket] %s", s)
}
