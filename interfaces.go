package intelhex

import (
	"io"
)

type RandomAccessReader interface {
	Read(location int) (byte, error)
}

type RandomAccessWriter interface {
	Write(location int, value byte) error
}

type RandomAccessReaderCloser interface {
	RandomAccessReader
	io.Closer
}

type RandomAccessWriterCloser interface {
	RandomAccessWriter
	io.Closer
}
