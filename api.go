package intelhex

import "io"

type ByteBlock struct {
	Address uint16
	Data    []byte
}

func Read(r io.Reader) ([]ByteBlock, error) {
	return nil, nil
}

func Write(w io.Writer, blocks []ByteBlock) error {

	for n := range blocks {
		err := WriteBlock(w, blocks[n])
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteBlock(w io.Writer, block ByteBlock) error {

	offset := 0
	for offset < len(block.Data) {
		newOffset, newAddr, err := writeDataLine(w, block.Data, block.Address, offset, 16)
		if err != nil {
			return err
		}
		offset = newOffset
		block.Address = newAddr
	}

	return nil
}

func WriteEOF(w io.Writer) error {
	return nil
}
