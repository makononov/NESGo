package cpu

import "fmt"

func (c *CPU) bit(val uint8) {
	test := val & c.a
	c.setZero(test)
	c.setNegative(val)
	c.overflow = (val&1<<6 != 0)
}

func (c *CPU) bmi(offset int8) {
	if c.negative {
		c.pc = uint16(int(c.pc) - int(offset))
	} else {
		c.pc += 2
	}
}

func (c *CPU) cld() {
	c.decimal = false
}

func (c *CPU) inc(address uint16) {
	value := c.readMem(address)
	value++
	c.setZero(value)
	c.setNegative(value)
	c.writeMem(address, value)
}

func (c *CPU) jmp(address uint16) {
	c.pc = address
}

func (c *CPU) jsr(address uint16) {
	lowByte := uint8((c.pc - 1) & 0xff)
	highByte := uint8((c.pc - 1) >> 8)
	c.stackPush(lowByte)
	c.stackPush(highByte)
	c.pc = address
}

func (c *CPU) lda(val uint8) {
	c.a = val
	c.setNegative(val)
	c.setZero(val)
}

func (c *CPU) ldx(val uint8) {
	c.x = val
	c.setNegative(val)
	c.setZero(val)
}

func (c *CPU) pha() {
	c.stackPush(c.a)
}

func (c *CPU) pla() {
	value := c.stackPop()
	c.a = value
	c.setNegative(value)
	c.setZero(value)
}

func (c *CPU) sbc(val uint8) {
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

func (c *CPU) sec() {
	c.carry = true
}

func (c *CPU) sei() {
	c.interruptDisable = true
}

func (c *CPU) sta(address uint16) {
	c.writeMem(address, c.a)
}

func (c *CPU) stx(address uint16) {
	c.writeMem(address, c.x)
}

func (c *CPU) sty(address uint16) {
	c.writeMem(address, c.y)
}
