package cpu

import "fmt"

// Operation is a function with an option addressing mode that executes the corresponding opcode
type Operation func(AddressingMode)

func (c *CPU) branchif(flag bool, offset int8) {
	if flag {
		newPC := Address(int(c.pc+2) + int(offset))
		if page(newPC) != page(c.pc) {
			c.cycleCount++
		}
		c.pc = newPC
		c.cycleCount++
	} else {
		c.pc += 2
	}
}

func (c *CPU) adc(address AddressingMode) {
	var carry uint8
	val := c.readMem(address())
	if c.carry {
		carry = 1
	}
	trueResult := uint16(c.a) + uint16(val) + uint16(carry)
	c.carry = (trueResult > 0xff)

	signBit := val >> 7
	c.overflow = (signBit == 1 && trueResult >= 0) || (signBit == 0 && trueResult < 0)

	c.a = val
	c.setZero(c.a)
	c.setNegative(c.a)
}

func (c *CPU) and(address AddressingMode) {
	val := c.readMem(address())
	c.a = c.a & val
	c.setZero(c.a)
	c.setNegative(c.a)
}

func (c *CPU) bcc(relative AddressingMode) {
	offset := int8(c.readMem(relative()))
	c.branchif(!c.carry, offset)
}

func (c *CPU) bcs(relative AddressingMode) {
	offset := int8(c.readMem(relative()))
	c.branchif(c.carry, offset)
}

func (c *CPU) beq(relative AddressingMode) {
	offset := int8(c.readMem(relative()))
	c.branchif(c.zero, offset)
}

func (c *CPU) bit(address AddressingMode) {
	val := c.readMem(address())
	test := val & c.a
	c.setZero(test)
	c.setNegative(val)
	c.overflow = (val&1<<6 != 0)
}

func (c *CPU) bmi(relative AddressingMode) {
	offset := int8(c.readMem(relative()))
	c.branchif(c.negative, offset)
}

func (c *CPU) bne(relative AddressingMode) {
	offset := int8(c.readMem(relative()))
	c.branchif(!c.zero, offset)
}

func (c *CPU) bpl(relative AddressingMode) {
	offset := int8(c.readMem(relative()))
	c.branchif(!c.negative, offset)
}

func (c *CPU) clc(_ AddressingMode) {
	c.carry = false
}

func (c *CPU) cld(_ AddressingMode) {
	c.decimal = false
}

func (c *CPU) cmp(address AddressingMode) {
	value := c.readMem(address())
	result := c.a - value
	c.carry = (c.a >= value)
	c.setZero(result)
	c.setNegative(result)
}

func (c *CPU) cpx(address AddressingMode) {
	value := c.readMem(address())
	result := c.x - value
	c.carry = (c.x >= value)
	c.setZero(result)
	c.setNegative(result)
}

func (c *CPU) cpy(address AddressingMode) {
	value := c.readMem(address())
	result := c.y - value
	c.carry = (c.y >= value)
	c.setZero(result)
	c.setNegative(result)
}

func (c *CPU) dec(address AddressingMode) {
	val := c.readMem(address())
	val--
	c.setZero(val)
	c.setNegative(val)
	c.writeMem(address(), val)
}

func (c *CPU) dex(_ AddressingMode) {
	c.x--
	c.setZero(c.x)
	c.setNegative(c.x)
}

func (c *CPU) dey(_ AddressingMode) {
	c.y--
	c.setZero(c.y)
	c.setNegative(c.y)
}

func (c *CPU) inc(address AddressingMode) {
	value := c.readMem(address())
	value++
	c.setZero(value)
	c.setNegative(value)
	c.writeMem(address(), value)
}

func (c *CPU) inx(_ AddressingMode) {
	c.x++
	c.setZero(c.x)
	c.setNegative(c.x)
}

func (c *CPU) iny(_ AddressingMode) {
	c.y++
	c.setZero(c.y)
	c.setNegative(c.y)
}

func (c *CPU) jmp(address AddressingMode) {
	c.pc = address()
}

func (c *CPU) jsr(address AddressingMode) {
	lowByte := uint8((c.pc + 2) & 0xff)
	highByte := uint8((c.pc + 2) >> 8)
	c.stackPush(lowByte)
	c.stackPush(highByte)
	c.pc = address()
}

func (c *CPU) lda(address AddressingMode) {
	c.a = c.readMem(address())
	c.setNegative(c.a)
	c.setZero(c.a)
}

func (c *CPU) ldx(address AddressingMode) {
	c.x = c.readMem(address())
	c.setNegative(c.x)
	c.setZero(c.x)
}

func (c *CPU) ldy(address AddressingMode) {
	c.y = c.readMem(address())
	c.setNegative(c.y)
	c.setZero(c.y)
}

func (c *CPU) lsr(address AddressingMode) {
	var value uint8
	if address == nil {
		value = c.a
	} else {
		value = c.readMem(address())
	}

	newVal := value >> 1
	c.carry = (value&1 != 0)
	c.setZero(newVal)
	c.setNegative(newVal)

	if address == nil {
		c.a = value
	} else {
		c.writeMem(address(), value)
	}
}

func (c *CPU) nop(_ AddressingMode) {}

func (c *CPU) ora(address AddressingMode) {
	val := c.readMem(address())
	c.a = c.a | val
	c.setZero(c.a)
	c.setNegative(c.a)
}

func (c *CPU) pha(_ AddressingMode) {
	c.stackPush(c.a)
}

func (c *CPU) pla(_ AddressingMode) {
	c.a = c.stackPop()
	c.setNegative(c.a)
	c.setZero(c.a)
}

func (c *CPU) rol(address AddressingMode) {
	var value uint8
	if address == nil {
		value = c.a
	} else {
		value = c.readMem(address())
	}

	newCarry := value >> 7
	newVal := value << 1
	if c.carry {
		newVal |= 1
	}
	c.carry = (newCarry == 1)
	c.setNegative(newVal)
	c.setZero(newVal)

	if address == nil {
		c.a = newVal
	} else {
		c.writeMem(address(), newVal)
	}
}

func (c *CPU) rts(_ AddressingMode) {
	highByte := uint16(c.stackPop())
	lowByte := uint16(c.stackPop())
	c.pc = Address(highByte<<8 | lowByte)
	c.pc++
}

func (c *CPU) sbc(address AddressingMode) {
	val := c.readMem(address())
	var t int8
	if c.decimal {
		t = binToBcd(c.a) - binToBcd(val)
		if !c.carry {
			t--
		}
		c.overflow = (t > 99 || t < 0)
	} else {
		t = int8(c.a - val)
		if !c.carry {
			t--
		}
		c.overflow = (c.a < val || (c.a == val) && !c.carry)
	}
	c.carry = (t >= 0)
	c.setNegative(uint8(t))
	c.setZero(uint8(t))
	fmt.Printf("SBC: 0x%x - 0x%x = 0x%x, C:%t N:%t V:%t Z:%t\n", c.a, val, uint8(t), c.carry, c.negative, c.overflow, c.zero)
	c.a = uint8(t)
}

func (c *CPU) sec(_ AddressingMode) {
	c.carry = true
}

func (c *CPU) sed(_ AddressingMode) {
	c.decimal = true
}

func (c *CPU) sei(_ AddressingMode) {
	c.interruptDisable = true
}

func (c *CPU) sta(address AddressingMode) {
	c.writeMem(address(), c.a)
}

func (c *CPU) stx(address AddressingMode) {
	c.writeMem(address(), c.x)
}

func (c *CPU) sty(address AddressingMode) {
	c.writeMem(address(), c.y)
}

func (c *CPU) tax(_ AddressingMode) {
	c.x = c.a
	c.setZero(c.x)
	c.setNegative(c.x)
}

func (c *CPU) tay(_ AddressingMode) {
	c.y = c.a
	c.setZero(c.y)
	c.setNegative(c.y)
}

func (c *CPU) tsx(_ AddressingMode) {
	c.x = c.sp
	c.setZero(c.x)
	c.setNegative(c.x)
}

func (c *CPU) txa(_ AddressingMode) {
	c.a = c.x
	c.setZero(c.a)
	c.setNegative(c.a)
}

func (c *CPU) txs(_ AddressingMode) {
	c.sp = c.x
}

func (c *CPU) tya(_ AddressingMode) {
	c.a = c.y
	c.setZero(c.a)
	c.setNegative(c.a)
}
