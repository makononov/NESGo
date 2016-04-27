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

	// status flags
	carry            bool
	zero             bool
	interruptDisable bool
	decimal          bool
	overflow         bool
	negative         bool

	// I/O Registers
	dmcraw   uint8
	dmcstart uint8
	dmclen   uint8

	// Communication busses
	cartridgeControlBus chan uint16
	ppuControlBus       chan uint16
	readWriteBus        chan int
	dataBus             chan uint8

	cycleCount int
}

// Init sets the CPU values to their initial power-up state.
func (c *CPU) Init(ppuControlBus chan uint16, cartridgeControlBus chan uint16, readWriteBus chan int, dataBus chan uint8) {
	c.ppuControlBus = ppuControlBus
	c.cartridgeControlBus = cartridgeControlBus
	c.readWriteBus = readWriteBus
	c.dataBus = dataBus
	c.ram = make([]byte, 2048)
	c.pc = 0

	c.interruptDisable = true
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
		fmt.Printf("0x%x", c.pc)
		opcode = c.readMem(c.pc)
		c.execute(opcode)
	}
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

func (c *CPU) writeMem(address uint16, val uint8) error {
	var controlBus chan uint16

	if address >= 0x8000 { // Cartridge ROM
		controlBus = c.cartridgeControlBus
	} else if address >= 0x2000 && address < 0x4000 {
		controlBus = c.ppuControlBus
	} else if address < 0x0800 { // Internal RAM
		c.ram[address] = val
		return nil
	} else {
		return fmt.Errorf("Attempt to write to an invalid memory location: 0x%x\n", address)
	}

	controlBus <- address
	c.readWriteBus <- 1 // write
	c.dataBus <- val

	return nil
}

func (c *CPU) execute(opcode uint8) {
	switch opcode {
	case 0x2c: // BIT absolute
		fmt.Printf(": BIT $%x (0x%x)\n", c.absolute(), c.readMem(c.absolute()))
		c.bit(c.readMem(c.absolute()))
		c.cycleCount += 3
		c.pc += 2
		break
	case 0x78: // SEI
		fmt.Printf(": SEI\n")
		c.sei()
		c.cycleCount += 2
		c.pc++
		break
	case 0x85: // STA zeropage
		fmt.Printf(": STA $%x\n", c.zeropage())
		c.sta(c.zeropage())
		c.cycleCount += 3
		c.pc += 2
		break
	case 0x8d: // STA abs
		fmt.Printf(": STA $%x\n", c.absolute())
		c.sta(c.absolute())
		c.cycleCount += 4
		c.pc += 3
		break
	case 0xa2: // LDX Immediate
		fmt.Printf(": LDX #%x\n", c.immediate())
		c.ldx(c.immediate())
		c.cycleCount += 2
		c.pc += 2
		break
	case 0xa9: // LDA Immediate
		fmt.Printf(": LDA #%x\n", c.immediate())
		c.lda(c.immediate())
		c.cycleCount += 2
		c.pc += 2
		break
	case 0xd8: // CLD
		fmt.Printf(": CLD\n")
		c.cld()
		c.cycleCount += 2
		c.pc++
		break
	default:
		panic(fmt.Errorf("Opcode not supported: 0x%x\n", opcode))
	}
}
