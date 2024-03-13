package utils

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

/*
creates a new PacketWriter (see also PacketReader)
The writer, writes packets to the stream. It assumes each call to write is an entire packet.
any start/stop or escape bytes in the payload will be escaped with the escape byte
*/
func NewPacketWriter(r io.Writer, start, escape, stop byte) (*PacketWriter, error) {
	err := check_valid(start, escape, stop)
	if err != nil {
		return nil, err
	}
	res := &PacketWriter{
		start:  start,
		stop:   stop,
		esc:    escape,
		writer: r,
	}
	return res, nil
}

type PacketWriter struct {
	start  byte
	stop   byte
	esc    byte
	writer io.Writer
	closed bool
	wrlock sync.Mutex
}

func (pr *PacketWriter) Close() error {
	pr.closed = true
	iw, castable := pr.writer.(io.WriteCloser)
	if castable {
		err := iw.Close()
		return err
	}
	return nil
}

// write a packet (using buf as the payload)
func (pr *PacketWriter) Write(buf []byte) (int, error) {
	if pr.closed {
		return 0, fmt.Errorf("write to closed packet writer")
	}
	if len(buf) > 8192 {
		return 0, fmt.Errorf("exceeded max packet size")
	}

	// TODO: convert to streaming for more efficiency
	pkt := bytes.Buffer{}
	pkt.Write([]byte{pr.start})
	for _, b := range buf {
		if b == pr.start || b == pr.stop || b == pr.esc {
			pkt.Write([]byte{pr.esc})
		}
		pkt.Write([]byte{b})
	}
	pkt.Write([]byte{pr.stop})
	pktbytes := pkt.Bytes()

	pr.wrlock.Lock() // start lock for atomic write for any one packet
	n, err := pr.writer.Write(pktbytes)
	pr.wrlock.Unlock() // end lock for atomic write for any one packet

	if err != nil {
		return n, err
	}
	return len(buf), nil
}
