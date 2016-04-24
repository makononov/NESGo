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

	// Communication busses
	cartridgeControlBus chan uint16
	readWriteBus        chan int
	dataBus             chan uint8

	cycleCount int
}

// Init sets the CPU values to their initial power-up state.
func (c *CPU) Init(cartridgeControlBus chan uint16, readWriteBus chan int, dataBus chan uint8) {
	c.cartridgeControlBus = cartridgeControlBus
	c.readWriteBus = readWriteBus
	c.dataBus = dataBus
	c.ram = make([]byte, 2048)
	c.p = 0x34
	c.pc = 0
}

// Run is the main function that processes through the PRG ROM.
func (c *CPU) Run() {
	fmt.Println("CPU spawned, getting initial PC...")
	lowByte := uint16(c.readMem(0xfffc))
	highByte := uint16(c.readMem(0xfffd))
	c.pc = highByte<<8 | lowByte

	fmt.Println("Beginning execution loop...")
	var opcode uint8
	for {
		opcode = c.readMem(c.pc)
		c.execute(opcode)
	}
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

func (c *CPU) readMem(address uint16) uint8 {
	var controlBus chan uint16

	if address >= 0x8000 {
		controlBus = c.cartridgeControlBus
	} else {
		panic(fmt.Errorf("Attempt to read invalid memory location: 0x%x\n", address))
	}

	controlBus <- address
	c.readWriteBus <- 0 // read
	return <-c.dataBus
}

func (c *CPU) writeMem(address uint16, val uint8) {
	var controlBus chan uint16

	if address >= 0x8000 { // Cartridge ROM
		controlBus = c.cartridgeControlBus
	} else if address >= 0x0200 && address < 0x0800 { // Internal RAM
		relativeAddress := address - 0x0200
		c.ram[relativeAddress] = val
		return
	} else {
		panic(fmt.Errorf("Attempt to write to an invalid memory location: 0x%x\n", address))
	}

	controlBus <- address
	c.readWriteBus <- 1 // write
}

func (c *CPU) execute(opcode uint8) {
	switch opcode {
	case 0x78: // SEI
		c.sei()
		c.cycleCount += 2
		c.pc += 1
		break
	case 0x8d: //STA abs
		c.sta(c.absolute())
		c.cycleCount += 4
		c.pc += 3
		break
	case 0xa9: //LDA Immediate
		c.lda(c.immediate())
		c.cycleCount += 2
		c.pc += 2
		break
	default:
		panic(fmt.Errorf("Opcode not supported: 0x%x\n", opcode))
	}
}
