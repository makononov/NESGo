package cpu

import "fmt"

// CPU emulates the 6502 processor
type CPU struct {
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
func (c *CPU) Init() {
	c.ram = make([]byte, 2048)
	c.p = 0x34
	c.pc = 0
}

// Run is the main function that processes through the PRG ROM.
func (c *CPU) Run(cartridgeControlBus chan uint16, readWriteBus chan int, dataBus chan uint8) {
	fmt.Println("CPU spawned, getting initial PC...")
	cartridgeControlBus <- 0xfffc
	readWriteBus <- 0 // read
	c.pc = c.pc | uint16(<-dataBus)<<2
	cartridgeControlBus <- 0xfffd
	readWriteBus <- 0
	c.pc = c.pc | uint16(<-dataBus)
	fmt.Printf("PC set to %x", c.pc)
}

func (c *CPU) carry() bool {
	return c.p&0x01 != 0
}

func (c *CPU) setCarry() {
	c.p = c.p | 0x01
}

func (c *CPU) zero() bool {
	return c.p&0x02 != 0
}

func (c *CPU) setZero() {
	c.p = c.p | 0x02
}

func (c *CPU) interruptDisable() bool {
	return c.p&0x04 != 0
}

func (c *CPU) setInterruptDisable() {
	c.p = c.p | 0x04
}

func (c *CPU) decimal() bool {
	return c.p&0x08 != 0
}

func (c *CPU) setDecimal() {
	c.p = c.p | 0x08
}

func (c *CPU) breakCommand() bool {
	return c.p&0x10 != 0
}

func (c *CPU) setBreakCommand() {
	c.p = c.p | 0x10
}

func (c *CPU) overflow() bool {
	return c.p&0x40 != 0
}

func (c *CPU) setOverflow() {
	c.p = c.p | 0x40
}

func (c *CPU) negative() bool {
	return c.p&0x80 != 0
}

func (c *CPU) setNegative() {
	c.p = c.p | 0x80
}
