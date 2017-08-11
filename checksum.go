package intelhex

type checksum struct {
	val byte
}

func (c *checksum) addBytes(bytes []byte) {
	for _, v := range bytes {
		c.val += v
	}
}

func (c *checksum) addByte(b byte) {
	c.val += b
}

func (c *checksum) addWord(w uint16) {
	c.addByte(byte(w))
	c.addByte(byte(w >> 8))
}

func (c *checksum) clear() {
	c.val = 0
}

func (c *checksum) value() byte {
	return -c.val
}
