package intelhex

import (
	"fmt"
	"io"
	"strconv"
)

func isContiguous(b1, b2 ByteBlock) bool {
	return (b1.Address+uint16(len(b1.Data)) == b2.Address) ||
		(b2.Address+uint16(len(b2.Data)) == b1.Address)
}

// assumes blocks are contiguous - order of supplied blocks not important
func joinByteBlocks(b1, b2 ByteBlock) ByteBlock {
	if b1.Address > b2.Address {
		b1, b2 = b2, b1
	}

	r := make([]byte, len(b1.Data)+len(b2.Data))
	copy(r, b1.Data)
	copy(r[len(b1.Data):], b2.Data)
	return ByteBlock{Address: b1.Address, Data: r}
}

func verifyCheckSum(data []byte) bool {
	chk := checksum{}
	chk.addBytes(data[:len(data)-1])

	return chk.value() == data[len(data)-1]
}

func processLineData(line string) (ByteBlock, error) {
	result := ByteBlock{}

	if len(line) < 11 {
		return result, fmt.Errorf("input line too short: %s", line)
	}

	if line[0] != ':' {
		return result, fmt.Errorf("input line does not start with colon: %s", line)
	}

	data, err := hexStrToBytes(line[1:])
	if err != nil {
		return result, err
	}

	if !verifyCheckSum(data) {
		return result, fmt.Errorf("input line bad checksum: %s", line)
	}

	dataLen := data[0]
	addr := uint16(data[1])<<8 + uint16(data[2])
	recType := data[3]

	if recType == 0 {
		result.Address = addr
		result.Data = data[4 : 4+dataLen]
	} else if recType == 1 {
		return result, io.EOF
	} else {
		return result, fmt.Errorf("Unsupported record type: %02x", recType)
	}

	return result, nil
}

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
	nextAddr = address + uint16(length)
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
