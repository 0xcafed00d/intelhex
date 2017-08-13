package intelhex

import (
	"bytes"
	"testing"

	"github.com/simulatedsimian/assert"
)

func TestCheckSum(t *testing.T) {
	assert := assert.Make(t)
	c := checksum{}

	assert(c.value()).Equal(byte(0x00))

	c.addByte(1)
	assert(c.value()).Equal(byte(0xff))

	c.clear()
	assert(c.value()).Equal(byte(0x00))

	c.addBytes([]byte{0x03, 0x00, 0x30, 0x00, 0x02, 0x33, 0x7A})
	assert(c.value()).Equal(byte(0x1e))

}

func TestStrToHex(t *testing.T) {
	assert := assert.Make(t)

	assert(hexStrToBytes("")).Equal([]byte{}, nil)
	assert(hexStrToBytes("1122ff")).Equal([]byte{0x11, 0x22, 0xff}, nil)
	assert(hexStrToBytes("11k2ff")).HasError()
	assert(hexStrToBytes("1122fff")).HasError()
}

func TestLineOutput(t *testing.T) {
	assert := assert.Make(t)

	buf := &bytes.Buffer{}
	writeDataLine(buf, []byte{0x02, 0x33, 0x7A}, 0x0030, 0, 16)
	assert(buf.String()).Equal(":0300300002337A1E\n")

	buf.Reset()
	writeDataLine(buf, []byte{0x02, 0x33, 0x7A}, 0x0030, 0, 3)
	assert(buf.String()).Equal(":0300300002337A1E\n")

	buf.Reset()
	writeDataLine(buf, []byte{0x02, 0x33, 0x7A}, 0x0030, 0, 1)
	assert(buf.String()).Equal(":0100300002CD\n")
}
