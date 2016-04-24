package cpu

func (c *CPU) lda(val uint8) {
	c.a = val
}

func (c *CPU) sei() {
	c.setInterruptDisable()
}

func (c *CPU) sta(address uint16) {
	c.writeMem(address, c.a)
}
