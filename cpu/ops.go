package cpu

func (c *CPU) bit(val uint8) {
	test := val & c.a
	c.zero = (test == 0)
}

func (c *CPU) cld() {
	c.decimal = false
}

func (c *CPU) lda(val uint8) {
	c.a = val
}

func (c *CPU) ldx(val uint8) {
	c.x = val
}

func (c *CPU) sei() {
	c.interruptDisable = true
}

func (c *CPU) sta(address uint16) {
	c.writeMem(address, c.a)
}
