package utils

import (
	"fmt"
	"io"
	"sync"
)

/*
This is part of the Packetizer toolset.

A PacketReader reads from any stream (specifically, an io.Reader). It reads from the stream until any one of the following conditions occur:

 1. the io.Reader signals io.EOF
 2. the io.Reader signals some other error
 3. a complete packet, starting with the start byte and ending with the stop byte has been received through the io.Reader
 4. the start of a packet has been read, but its size exceeds the buffer size (specifically: the of the array passed to PacketReader.Read([]byte))

A packet is defined as the data between start and stop byte. As a consequence, Any start and stop bytes in the payload must be
escaped with the escape byte

The data returned by the packet reader is guaranteed to be equal to the payload send by the PacketWriter. In other words:
any escaping and wrapping and unescaping and unwrapping is handled by the reader and writer.
see also [NewPacketWriter]
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
	stop_lock    sync.Mutex
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
	pr.stopped = true
	pr.stop_lock.Lock()
	if !pr.stopped {
		close(pr.packet_chan)
		pr.stopped = true
	}
	pr.stop_lock.Unlock()
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
		if pr.stopped {
			break
		}
		n, err := pr.reader.Read(buf)
		if pr.stopped {
			break
		}
		if n > 0 {
			payload := buf[:n]
			pr.debugf("Read %d bytes:\n%s\n", n, Hexdump("payload: ", payload))
			for _, b := range payload {
				if stp.AddByte(b) {
					pkt := stp.ReadPacket()
					pr.debugf("%s\n", Hexdump("packet:", pkt))
					pr.stop_lock.Lock()
					if pr.stopped {
						pr.stop_lock.Unlock()
						break
					}
					pr.packet_chan <- &readerEvent{payload: pkt}
					pr.stop_lock.Unlock()
				}
			}
		}

		if err != nil {
			pr.stop_lock.Lock()
			if pr.stopped {
				pr.stop_lock.Unlock()
				break
			}
			pr.packet_chan <- &readerEvent{err: err}
			pr.stop_lock.Unlock()
			break
		}
	}
	pr.debugf("read_loop stopped\n")
}

func (stp *PacketReader) debugf(format string, args ...interface{}) {
	if !stp.debug {
		return
	}
	s := fmt.Sprintf(format, args...)
	fmt.Printf("[packetreader] %s", s)
}
