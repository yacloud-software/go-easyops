package utils

import (
	"fmt"
	"io"
)

/*
creates a new PacketReader.
The reader, reads from a stream (an io.Reader). Reads from the PacketReader block until a complete packet is received (or EOF).
each read is guaranteed to return either complete packet in the buf or error if the packet won't fit into the buf.
a packet is started with the start byte, ends with the stop byte.
any start/stop or escape bytes in the payload must be escaped with the escape byte
it automatically resynchronises on lost data. This means, it is suitable for lossy transmission layers, such as UDP or serial ports.
the lost data is not recovered (it is lost), but it will continue to parse the stream and resynchronise to start/stop bytes.
any packet that is returned to the reader, is already un-escaped.
*/
func NewPacketReader(r io.Reader, start, escape, stop byte) (*PacketReader, error) {
	err := check_valid(start, escape, stop)
	if err != nil {
		return nil, err
	}
	res := &PacketReader{
		start:       start,
		stop:        stop,
		esc:         escape,
		debug:       *debug_packetizer,
		reader:      r,
		packet_chan: make(chan *readerEvent, 10),
	}
	go res.packet_reader_loop()
	return res, nil
}

type PacketReader struct {
	start        byte
	stop         byte
	esc          byte
	reader       io.Reader
	debug        bool
	packet_chan  chan *readerEvent
	reader_error error // any error the io.reader has returned
	stopped      bool
}

type readerEvent struct {
	payload []byte
	err     error
}

// this will return the number of bytes or an error. it is guaranteed to either return an error OR a non-zero number of bytes, but never both
func (pr *PacketReader) Read(buf []byte) (int, error) {
	if pr.reader_error != nil {
		return 0, pr.reader_error
	}

	packet := <-pr.packet_chan
	if packet.err != nil {
		pr.reader_error = packet.err
		return 0, pr.reader_error
	}
	n := len(packet.payload)
	if n > len(buf) {
		return 0, fmt.Errorf("packetreader buf too small. got %d bytes, but need at least %d bytes", len(buf), n)
	}
	for i := 0; i < n; i++ {
		buf[i] = packet.payload[i]
	}
	pr.debugf("Got packet %d bytes\n", n)
	return n, nil
}

// closes underlying reader as well AND aborts a current read
func (pr *PacketReader) Close() {
	if pr.stopped {
		return
	}
	rc, cast := pr.reader.(io.ReadCloser)
	if cast {
		rc.Close()
	}
	pr.packet_chan <- &readerEvent{err: io.EOF}
	close(pr.packet_chan)
}

// copy from io.reader to packets
func (pr *PacketReader) packet_reader_loop() {
	buf := make([]byte, 200)
	stp, err := NewStreamToPacket(pr.start, pr.esc, pr.stop)
	if err != nil {
		pr.reader_error = err
		pr.packet_chan <- &readerEvent{err: err}
		return
	}

	for {
		n, err := pr.reader.Read(buf)
		if n > 0 {
			payload := buf[:n]
			pr.debugf("Read %d bytes:\n%s\n", n, Hexdump("payload: ", payload))
			for _, b := range payload {
				if stp.AddByte(b) {
					pkt := stp.ReadPacket()
					pr.debugf("%s\n", Hexdump("packet:", pkt))
					pr.packet_chan <- &readerEvent{payload: pkt}
				}
			}
		}

		if err != nil {
			pr.packet_chan <- &readerEvent{err: err}
			break
		}

	}
}

func (stp *PacketReader) debugf(format string, args ...interface{}) {
	if !stp.debug {
		return
	}
	s := fmt.Sprintf(format, args...)
	fmt.Printf("[packetreader] %s", s)
}
