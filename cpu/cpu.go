package cpu

import (
	"fmt"
	"math"
)

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
	c.sp = 0xff

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
	executed := 0
	for {
		executed++
		fmt.Printf("%d 0x%x A:0x%x X:0x%x Y:0x%x SP:0x%x", executed, c.pc, c.a, c.x, c.y, c.sp)
		opcode = c.readMem(c.pc)
		c.execute(opcode)
	}
}

func (c *CPU) readMem(address uint16) uint8 {
	var controlBus chan uint16

	if address >= 0x8000 {
		controlBus = c.cartridgeControlBus
	} else if address >= 0x2000 && address < 0x4000 {
		controlBus = c.ppuControlBus
	} else if address < 0x0800 {
		return c.ram[address]
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
	case 0x10: // BPL
		fmt.Printf(": BPL $%x\n", c.relative())
		prevPC := c.pc
		c.bpl(c.relative())
		if c.pc == prevPC {
			// Branch was not taken
			c.cycleCount += 2
		} else {
			c.cycleCount++

			if page(prevPC) != page(c.pc) {
				// Crossed page boundary
				c.cycleCount++
			}
		}
		break
	case 0x18: // CLC
		fmt.Printf(": CLC\n")
		c.clc()
		c.cycleCount += 2
		c.pc++
		break
	case 0x20: // JSR absolute
		fmt.Printf(": JSR $%x\n", c.absolute())
		c.jsr(c.absolute())
		c.cycleCount += 6
		break
	case 0x2c: // BIT absolute
		fmt.Printf(": BIT $%x (0x%x)", c.absolute(), c.readMem(c.absolute()))
		c.bit(c.readMem(c.absolute()))
		c.cycleCount += 3
		c.pc += 2
		fmt.Printf(" Z:%t S:%t V:%t\n", c.zero, c.negative, c.overflow)
		break
	case 0x30: // BMI
		fmt.Printf(": BMI $%x\n", c.relative())
		prevPC := c.pc
		c.bmi(c.relative())
		if c.pc == prevPC {
			// Branch was not taken
			c.cycleCount += 2
		} else {
			c.cycleCount++

			if page(prevPC) != page(c.pc) {
				// Crossed page boundary
				c.cycleCount++
			}
		}
		break
	case 0x38: // SEC
		fmt.Printf(": SEC\n")
		c.sec()
		c.cycleCount += 2
		c.pc++
		break
	case 0x48: // PHA
		fmt.Printf(": PHA\n")
		c.pha()
		c.cycleCount += 3
		c.pc++
		break
	case 0x4c: // JMP absolute
		fmt.Printf(": JMP $%x\n", c.absolute())
		c.jmp(c.absolute())
		c.cycleCount += 3
		break
	case 0x68: // PLA
		fmt.Printf(": PLA\n")
		c.pla()
		c.cycleCount += 4
		c.pc++
		break
	case 0x78: // SEI
		fmt.Printf(": SEI\n")
		c.sei()
		c.cycleCount += 2
		c.pc++
		break
	case 0x84: // STY zeropage
		fmt.Printf(": STY $%x\n", c.zeropage())
		c.sty(c.zeropage())
		c.cycleCount += 3
		c.pc += 2
		break
	case 0x85: // STA zeropage
		fmt.Printf(": STA $%x\n", c.zeropage())
		c.sta(c.zeropage())
		c.cycleCount += 3
		c.pc += 2
		break
	case 0x86: // STX zeropage
		fmt.Printf(": STX $%x\n", c.zeropage())
		c.stx(c.zeropage())
		c.cycleCount += 3
		c.pc += 2
		break
	case 0x88: // DEY
		fmt.Printf(": DEY\n")
		c.dey()
		c.cycleCount += 2
		c.pc++
		break
	case 0x8a: // TXA
		fmt.Printf(": TXA\n")
		c.txa()
		c.cycleCount += 2
		c.pc++
	case 0x8d: // STA abs
		fmt.Printf(": STA $%x\n", c.absolute())
		c.sta(c.absolute())
		c.cycleCount += 4
		c.pc += 3
		break
	case 0x91: // STA (Indirect),Y
		fmt.Printf(": STA ($%x), Y\n", c.immediate())
		c.sta(c.indirectY())
		c.cycleCount += 6
		c.pc += 2
		break
	case 0x9a: // TXS
		fmt.Printf(": TXS\n")
		c.txs()
		c.cycleCount += 2
		c.pc++
	case 0xa0: // LDY Immediate
		fmt.Printf(": LDY #%x\n", c.immediate())
		c.ldy(c.immediate())
		c.cycleCount += 2
		c.pc += 2
	case 0xa2: // LDX Immediate
		fmt.Printf(": LDX #%x\n", c.immediate())
		c.ldx(c.immediate())
		c.cycleCount += 2
		c.pc += 2
		break
	case 0xa5: // LDA Zeropage
		fmt.Printf(": LDA $%x\n", c.zeropage())
		c.lda(c.readMem(c.zeropage()))
		c.cycleCount += 3
		c.pc += 2
		break
	case 0xa9: // LDA Immediate
		fmt.Printf(": LDA #%x\n", c.immediate())
		c.lda(c.immediate())
		c.cycleCount += 2
		c.pc += 2
		break
	case 0xaa: // TAX
		fmt.Printf(": TAX\n")
		c.tax()
		c.cycleCount += 2
		c.pc++
		break
	case 0xad: // LDA Absolute
		fmt.Printf(": LDA $%x\n", c.absolute())
		c.lda(c.readMem(c.absolute()))
		c.cycleCount += 4
		c.pc += 3
	case 0xd8: // CLD
		fmt.Printf(": CLD\n")
		c.cld()
		c.cycleCount += 2
		c.pc++
		break
	case 0xe5: // SBC zeropage
		fmt.Printf(": SBC $%x\n", c.zeropage())
		c.sbc(c.readMem(c.zeropage()))
		c.cycleCount += 3
		c.pc += 2
		break
	case 0xe6: // INC zeropage
		fmt.Printf(": INC $%x\n", c.zeropage())
		c.inc(c.zeropage())
		c.cycleCount += 5
		c.pc += 2
		break
	default:
		panic(fmt.Errorf("Opcode not supported: 0x%x\n", opcode))
	}
}

func (c *CPU) stackPush(value uint8) {
	c.writeMem(0x0100+uint16(c.sp), value)
	c.sp--
}

func (c *CPU) stackPop() uint8 {
	c.sp++
	return c.readMem(0x0100 + uint16(c.sp))
}

func (c *CPU) setNegative(value uint8) {
	c.negative = (value&1<<7 != 0)
}

func (c *CPU) setZero(value uint8) {
	c.zero = (value == 0)
}

func binToBcd(val uint8) int8 {
	ones := int8(val & 0x0f)
	tens := int8(val & 0xf0 / 0x10)
	return (tens * 10) + ones
}

func page(address uint16) int {
	return int(math.Floor(float64(address) / 0x4000))
}
