package intelhex

import (
	"fmt"
	"io"
)

func writeDataLine(w io.Writer, data []byte, address uint16, offset, maxlen int) (nextOffset int, nextAddr uint16, err error) {
	c := checksum{}

	length := maxlen
	if length+offset > len(data) {
		length = len(data) - offset
	}

	c.addByte(byte(length))
	c.addWord(address)

	_, err = fmt.Fprintf(w, ":%02x%04x00", length, address)
	if err != nil {
		return
	}

	for n := 0; n < length; n++ {
		b := data[offset+n]
		c.addByte(b)
		_, err = fmt.Fprintf(w, "%02x", b)
		if err != nil {
			return
		}
	}

	_, err = fmt.Fprintf(w, "%02x", c.value())
	if err != nil {
		return
	}

	nextOffset = offset + length
	nextAddr += uint16(length)
	return
}
