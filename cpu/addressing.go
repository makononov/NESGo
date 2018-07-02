package cpu

// Address represents a 16-bit memory address
type Address uint16

// An AddressingMode returns an address that can be used when executing an operation
type AddressingMode func() Address

func (c *CPU) implied() Address { return 0 }

func (c *CPU) immediate() Address {
	return c.pc + 1
}

func (c *CPU) absolute() Address {
	return Address(c.readBytes(c.pc + 1))
}

func (c *CPU) absoluteX() Address {
	baseAddress := c.absolute()
	finalAddress := Address(uint16(baseAddress) + uint16(c.x))
	if page(baseAddress) != page(finalAddress) {
		c.cycleCount++
	}
	return finalAddress
}

func (c *CPU) absoluteY() Address {
	baseAddress := c.absolute()
	finalAddress := Address(uint16(baseAddress) + uint16(c.y))
	if page(baseAddress) != page(finalAddress) {
		c.cycleCount++
	}
	return finalAddress
}

func (c *CPU) zeropage() Address {
	return Address(c.readMem(c.pc + 1))
}

func (c *CPU) zeropageX() Address {
	baseAddress := c.zeropage()
	return Address(uint8(baseAddress) + c.x)
}

func (c *CPU) zeropageY() Address {
	baseAddress := c.zeropage()
	return Address(uint8(baseAddress) + c.y)
}

func (c *CPU) relative() Address {
	return c.pc + 1
}

func (c *CPU) indirect() Address {
	baseAddress := c.absolute()
	return Address(c.readBytes(baseAddress))
}

// indirectX is Indexed Indirect addressing using the X register
func (c *CPU) indirectX() Address {
	baseAddress := c.readMem(c.pc + 1)
	address := Address(baseAddress + c.x)
	return Address(c.readBytes(address))
}

// indirectY is Indirect Indexed addressing using the Y register
func (c *CPU) indirectY() Address {
	baseAddress := Address(c.readBytes(c.zeropage()))
	finalAddress := Address(uint16(baseAddress) + uint16(c.y))
	if page(baseAddress) != page(finalAddress) {
		c.cycleCount++
	}
	return finalAddress
}
