package utils

import (
	"bytes"
	"io"
	"math/rand"
	"testing"
)

const (
	RANDSOURCE  = 1233
	MAX_SIZE    = 7000
	MAX_PACKETS = 10000
)

func TestSelectedPackets(t *testing.T) {
	packets := [][]byte{
		[]byte("abc"),
		[]byte("abcd"),
		[]byte("abcde"),
		[]byte("ab{cde"),
		[]byte("ab}cde"),
		[]byte("ab\\}cde"),
		[]byte("ab\\\\cde"),
		[]byte("ab\\{cde"),
	}
	test_packets(t, packets)
}

func TestRandomDataAndRandomSize(t *testing.T) {
	r := rand.New(rand.NewSource(RANDSOURCE))
	var packets [][]byte
	for i := 0; i < MAX_PACKETS; i++ {
		packlen := r.Intn(MAX_SIZE)
		packet := make([]byte, packlen)
		for i, _ := range packet {
			b := byte(r.Intn(255))
			packet[i] = b
		}
		packets = append(packets, packet) // store for later comparison
	}
	test_packets(t, packets)
}
func test_packets(t *testing.T, packets [][]byte) {
	b := &bytes.Buffer{}
	pw, _ := NewPacketWriter(b, '{', '\\', '}')
	for _, packet := range packets {
		_, err := pw.Write(packet)
		if err != nil {
			t.Fail()
			t.Logf("failed to write: %s\n", err)
			return
		}
	}
	err := pw.Close()
	if err != nil {
		t.Fail()
		t.Logf("failed to close: %s\n", err)
	}
	bb := b.Bytes()
	if len(bb) == 0 {
		t.Fail()
		t.Logf("empty buffer after writing packets")
		return
	}
	w := bytes.NewBuffer(bb)
	compare(t, w, packets)
}

// compare what the writer gives us with "expected"
func compare(t *testing.T, w io.Reader, expected [][]byte) {
	pr, _ := NewPacketReader(w, '{', '\\', '}')
	i := 0
	buf := make([]byte, 8192)
	for {
		n, err := pr.Read(buf)
		if err != nil {
			if err == io.EOF {
				//t.Logf("EOF on reader\n")
				break
			}
			t.Fail()
			t.Logf("Unexpected read error: %s\n", err)
			return
		}

		if i > len(expected) {
			t.Fail()
			t.Logf("read too many packets")
			return
		}
		packet := buf[:n]
		mp := expected[i]
		if !bytes.Equal(packet, mp) {
			t.Fail()
			t.Logf("packet mismatch (position %d). output in /tmp/packetizer_read.bin and /tmp/packetizer_expected.bin", i)
			write_packets(packet, mp)
			return
		}
		i++

	}
	if i != len(expected) {
		t.Fail()
		t.Logf("read too few packets. expected %d, read %d packets\n", len(expected), i)
	}
}

func write_packets(a, b []byte) {
	WriteFile("/tmp/packetizer_read.bin", a)
	WriteFile("/tmp/packetizer_expected.bin", b)
	s := Hexdump("", a)
	WriteFile("/tmp/packetizer_read.txt", []byte(s))
	s = Hexdump("", b)
	WriteFile("/tmp/packetizer_expected.txt", []byte(s))

}
