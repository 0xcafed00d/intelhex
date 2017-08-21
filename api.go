package intelhex

import (
	"bufio"
	"io"
	"strings"
)

type ByteBlock struct {
	Address uint16
	Data    []byte
}

func Read(r io.Reader) ([]ByteBlock, error) {
	blocks := []ByteBlock{}
	br := bufio.NewReader(r)

	for {
		line, err := br.ReadString('\n')
		line = strings.TrimSpace(line)

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		block, err := processLineData(line)
		if err != nil {
			return nil, err
		}
		if len(blocks) == 0 || !isContiguous(block, blocks[len(blocks)-1]) {
			blocks = append(blocks, block)
		} else {
			blocks[len(blocks)-1] = joinByteBlocks(block, blocks[len(blocks)-1])
		}
	}
	return blocks, nil
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
