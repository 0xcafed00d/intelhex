package intelhex

import (
	"fmt"
	"io"
	"strconv"
)

func writeDataLine(w io.Writer, data []byte, address uint16, offset, maxlen int) (nextOffset int, nextAddr uint16, err error) {
	chk := checksum{}

	length := maxlen
	if length+offset > len(data) {
		length = len(data) - offset
	}

	chk.addByte(byte(length))
	chk.addWord(address)

	_, err = fmt.Fprintf(w, ":%02X%04X00", length, address)
	if err != nil {
		return

	}

	for n := 0; n < length; n++ {
		b := data[offset+n]
		chk.addByte(b)
		_, err = fmt.Fprintf(w, "%02X", b)
		if err != nil {
			return
		}
	}

	_, err = fmt.Fprintf(w, "%02X\n", chk.value())
	if err != nil {
		return
	}

	nextOffset = offset + length
	nextAddr += uint16(length)
	return
}

func writeEOFLine(w io.Writer) error {
	_, err := fmt.Fprintln(w, ":00000001FF")
	return err
}

func hexStrToBytes(hex string) ([]byte, error) {
	result := []byte{}

	for len(hex) >= 2 {
		val, err := strconv.ParseUint(hex[0:2], 16, 8)
		if err != nil {
			return nil, err
		}
		result = append(result, byte(val))
		hex = hex[2:]
	}

	if len(hex) > 0 {
		return nil, fmt.Errorf("Uneven number of characters supplied")
	}

	return result, nil
}
