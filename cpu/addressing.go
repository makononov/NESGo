package cpu

func (c *CPU) immediate() uint8 {
	return c.readMem(c.pc + 1)
}

func (c *CPU) absolute() uint16 {
	lowByte := uint16(c.readMem(c.pc + 1))
	highByte := uint16(c.readMem(c.pc + 2))
	return (highByte<<8 | lowByte)
}

func (c *CPU) zeropage() uint16 {
	lowByte := uint16(c.readMem(c.pc + 1))
	return lowByte
}

func (c *CPU) relative() int8 {
	return int8(c.readMem(c.pc + 1))
}

func (c *CPU) indirectY() uint16 {
	lowByte := c.zeropage()
	return lowByte + uint16(c.y)
}
