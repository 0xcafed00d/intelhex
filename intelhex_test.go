package intelhex

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/simulatedsimian/assert"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

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

var str1 = `line1
line2
line3
`

func TestLineReader(t *testing.T) {
	r := strings.NewReader(str1)
	br := bufio.NewReader(r)

	for {
		line, err := br.ReadString('\n')
		line = strings.TrimSpace(line)
		fmt.Println(">>>:", string(line), ":<<<")
		if err == io.EOF {
			fmt.Println("EOF")
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func TestJoinBlocks(t *testing.T) {
	assert := assert.Make(t)

	b1 := ByteBlock{0, []byte{0, 1, 2, 3}}
	b2 := ByteBlock{4, []byte{4, 5, 6, 7}}
	b3 := ByteBlock{8, []byte{8, 9, 10, 11}}
	b4 := ByteBlock{0, []byte{0, 1, 2, 3, 4, 5, 6, 7}}

	assert(isContiguous(b1, b2)).Equal(true)
	assert(isContiguous(b2, b1)).Equal(true)
	assert(isContiguous(b1, b3)).Equal(false)
	assert(isContiguous(b3, b1)).Equal(false)

	assert(joinByteBlocks(b1, b2)).Equal(b4)
	assert(joinByteBlocks(b2, b1)).Equal(b4)
}

func TestProcessLine(t *testing.T) {
	assert := assert.Make(t)

	assert(processLineData(":0300300002337A1E")).NoError()
	assert(processLineData("0300300002337A1E")).HasError()
	assert(processLineData(":0300300002337A1F")).HasError()
	assert(processLineData(":030030")).HasError()
	assert(processLineData(":0300300002337pA1F")).HasError()

	assert(processLineData(":0300300002337A1E")).
		Equal(ByteBlock{0x30, []byte{0x02, 0x33, 0x7a}}, nil)

	assert(processLineData(":00000001FF")).
		Equal(ByteBlock{}, io.EOF)
}

func createRandomByteBlock() ByteBlock {
	bb := ByteBlock{}
	size := rand.Intn(64) + 20

	bb.Address = uint16(rand.Intn(0x7fff))
	bb.Data = make([]byte, size)
	for i := range bb.Data {
		bb.Data[i] = byte(rand.Intn(255))
	}

	return bb
}

func TestReadFile(t *testing.T) {
	assert := assert.Make(t)
	buf := bytes.Buffer{}

	bb := createRandomByteBlock()

	WriteBlock(&buf, bb)
	WriteEOF(&buf)
	fmt.Println(string(buf.Bytes()))

	assert(Read(&buf)).Equal(bb, nil)

}
