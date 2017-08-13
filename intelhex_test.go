package intelhex

import "testing"
import "os"
import "github.com/simulatedsimian/assert"

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

func TestStrToHex(t *testing.T) {
	assert := assert.Make(t)

	assert(hexStrToBytes("")).Equal([]byte{}, nil)
	assert(hexStrToBytes("1122ff")).Equal([]byte{0x11, 0x22, 0xff}, nil)
	assert(hexStrToBytes("11k2ff")).HasError()
	assert(hexStrToBytes("1122fff")).HasError()
}

func TestLineOutput(t *testing.T) {
	writeDataLine(os.Stdout, []byte{0x02, 0x33, 0x7A}, 0x0030, 0, 16)
}
