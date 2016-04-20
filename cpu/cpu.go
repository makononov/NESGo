package cpu

type cpu struct {
	mem []int

	// Registers
	pc int // Program Counter
	sp int // Stack Pointer
	a  int // Accumulator
	x  int // X-addressing register
	y  int // Y-addressing register
	p  int // Processor Status
}

// Run is the main function that processes through the PRG ROM.
func Run(dataBus chan int, controlBus chan int, addressBus chan int) {

}

func (c *cpu) carry() bool {
	return c.p&0x01 != 0
}

func (c *cpu) setCarry() {
	c.p = c.p | 0x01
}

func (c *cpu) zero() bool {
	return c.p&0x02 != 0
}

func (c *cpu) setZero() {
	c.p = c.p | 0x02
}

func (c *cpu) interruptDisable() bool {
	return c.p&0x04 != 0
}

func (c *cpu) setInterruptDisable() {
	c.p = c.p | 0x04
}

func (c *cpu) decimal() bool {
	return c.p&0x08 != 0
}

func (c *cpu) setDecimal() {
	c.p = c.p | 0x08
}

func (c *cpu) breakCommand() bool {
	return c.p&0x10 != 0
}

func (c *cpu) setBreakCommand() {
	c.p = c.p | 0x10
}

func (c *cpu) overflow() bool {
	return c.p&0x40 != 0
}

func (c *cpu) setOverflow() {
	c.p = c.p | 0x40
}

func (c *cpu) negative() bool {
	return c.p&0x80 != 0
}

func (c *cpu) setNegative() {
	c.p = c.p | 0x80
}
