package intelhex

import "testing"
import "os"

func TestCheckSum(t *testing.T) {
	c := checksum{}

	if c.value() != 0x00 {
		t.Fail()
	}

	c.addByte(1)
	if c.value() != 0xff {
		t.Fail()
	}

	c.clear()
	if c.value() != 0x00 {
		t.Fail()
	}

	c.addBytes([]byte{0x03, 0x00, 0x30, 0x00, 0x02, 0x33, 0x7A})
	if c.value() != 0x1e {
		t.Fail()
	}
}

func TestLineOutput(t *testing.T) {
	writeDataLine(os.Stdout, []byte{0x02, 0x33, 0x7A}, 0x0030, 0, 16)
}
