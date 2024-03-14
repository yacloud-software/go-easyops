package utils

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

/*
This is part of the Packetizer toolset.

A PacketWriter writes arbitrary data to an io.Writer. Each call to [PacketWriter.Write] is assumed to be one packet.
The PacketWriter sends extra data to the io.Writer to allow a [PacketReader] to identify and reassemble each packet.

The algorithm to send any one packet can be summarised like so:

1. Send one extra byte, the start- byte
2. Send the payload upto, but not including any start, stop or escape bytes in the payload
3. for any start, stop or escape byte, insert an extra escape byte
4. repeat send until all bytes are send
5. send one extra byte, the escape-byte

This algorithm is sufficient to send any data, 8-bit clean, across any io.Writer and reassemble it on the receiving side ot he
stream.
The canonical implementation to reassemble the packets is [NewPacketReader].
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
