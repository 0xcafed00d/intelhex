package intelhex

import "io"

type ByteBlock struct {
	Address uint16
	Data    []byte
}

func Read(r io.Reader) ([]ByteBlock, error) {
	return nil, nil
}

func Write(r io.Writer, blocks []ByteBlock) error {

	return nil
}

func WriteBlock(r io.Writer, blocks ByteBlock) error {

	return nil
}

func WriteEOF(r io.Writer) error {
	return nil
}
