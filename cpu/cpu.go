package cpu

type cpu struct {
	ram []byte

	// special Registers
	pc uint16 // Program Counter
	sp uint8  // Stack Pointer
	a  uint8  // Accumulator
	x  uint8  // X-addressing register
	y  uint8  // Y-addressing register
	p  uint8  // Processor Status

	// I/O Registers
	dmcraw   uint8
	dmcstart uint8
	dmclen   uint8
}

// Init sets the CPU values to their initial power-up state.
func (c *cpu) Init() {
	c.ram = make([]byte, 2048)
	c.p = 0x34
}

// Run is the main function that processes through the PRG ROM.
func (c *cpu) Run(clock chan bool, dataBus chan int, controlBus chan int, addressBus chan int) {
	exit := <-clock

	if exit {
		return
	}
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
