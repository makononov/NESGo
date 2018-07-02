package cpu

import (
	"fmt"
	"log"
	"math"
)

type op struct {
	f Operation
	a AddressingMode
	c int
	s uint16
}

// CPU emulates the 6502 processor
type CPU struct {
	ram []byte

	// special Registers
	pc Address // Program Counter
	sp uint8   // Stack Pointer
	a  uint8   // Accumulator
	x  uint8   // X-addressing register
	y  uint8   // Y-addressing register

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
	vblankBus           chan bool

	vblank     bool
	cycleCount int

	operations map[uint8]op
}

// Init sets the CPU values to their initial power-up state.
func (c *CPU) Init(ppuControlBus chan uint16, cartridgeControlBus chan uint16, readWriteBus chan int, dataBus chan uint8, vblankBus chan bool) {
	c.ppuControlBus = ppuControlBus
	c.cartridgeControlBus = cartridgeControlBus
	c.readWriteBus = readWriteBus
	c.dataBus = dataBus
	c.vblankBus = vblankBus
	c.vblank = false
	c.ram = make([]byte, 2048)
	c.pc = 0
	c.sp = 0xff

	c.interruptDisable = true

	c.operations = map[uint8]op{
		0x01: {f: c.ora, a: c.indirectX, c: 6, s: 2},
		0x05: {f: c.ora, a: c.zeropage, c: 3, s: 2},
		0x09: {f: c.ora, a: c.immediate, c: 2, s: 2},
		0x0d: {f: c.ora, a: c.absolute, c: 4, s: 3},
		0x10: {f: c.bpl, a: c.relative, c: 1, s: 0},
		0x11: {f: c.ora, a: c.indirectY, c: 5, s: 2},
		0x15: {f: c.ora, a: c.zeropageX, c: 4, s: 2},
		0x18: {f: c.clc, a: c.implied, c: 2, s: 1},
		0x19: {f: c.ora, a: c.absoluteY, c: 4, s: 3},
		0x1d: {f: c.ora, a: c.absoluteX, c: 4, s: 3},
		0x20: {f: c.jsr, a: c.absolute, c: 6, s: 0},
		0x21: {f: c.and, a: c.indirectX, c: 6, s: 2},
		0x25: {f: c.and, a: c.zeropage, c: 3, s: 2},
		0x26: {f: c.rol, a: c.zeropage, c: 5, s: 2},
		0x29: {f: c.and, a: c.immediate, c: 2, s: 2},
		0x2a: {f: c.rol, a: nil, c: 2, s: 1},
		0x2c: {f: c.bit, a: c.absolute, c: 3, s: 2},
		0x2d: {f: c.and, a: c.absolute, c: 4, s: 3},
		0x2e: {f: c.rol, a: c.absolute, c: 6, s: 3},
		0x30: {f: c.bmi, a: c.relative, c: 1, s: 0},
		0x31: {f: c.and, a: c.indirectY, c: 5, s: 2},
		0x35: {f: c.and, a: c.zeropageX, c: 4, s: 2},
		0x36: {f: c.rol, a: c.zeropageX, c: 6, s: 2},
		0x38: {f: c.sec, a: c.implied, c: 2, s: 1},
		0x39: {f: c.and, a: c.absoluteY, c: 4, s: 3},
		0x3d: {f: c.and, a: c.absoluteX, c: 4, s: 3},
		0x3e: {f: c.rol, a: c.absoluteX, c: 7, s: 3},
		0x46: {f: c.lsr, a: c.zeropage, c: 5, s: 2},
		0x48: {f: c.pha, a: c.implied, c: 3, s: 1},
		0x4a: {f: c.lsr, a: nil, c: 2, s: 1},
		0x4c: {f: c.jmp, a: c.absolute, c: 3, s: 0},
		0x4e: {f: c.lsr, a: c.absolute, c: 6, s: 3},
		0x56: {f: c.lsr, a: c.zeropageX, c: 6, s: 2},
		0x5e: {f: c.lsr, a: c.absoluteX, c: 7, s: 3},
		0x60: {f: c.rts, a: c.implied, c: 6, s: 0},
		0x65: {f: c.adc, a: c.zeropage, c: 3, s: 2},
		0x68: {f: c.pla, a: c.implied, c: 4, s: 1},
		0x78: {f: c.sei, a: c.implied, c: 2, s: 1},
		0x84: {f: c.sty, a: c.zeropage, c: 3, s: 2},
		0x85: {f: c.sta, a: c.zeropage, c: 3, s: 2},
		0x86: {f: c.stx, a: c.zeropage, c: 3, s: 2},
		0x88: {f: c.dey, a: c.implied, c: 2, s: 1},
		0x8a: {f: c.txa, a: c.implied, c: 2, s: 1},
		0x8c: {f: c.sty, a: c.absolute, c: 4, s: 3},
		0x8d: {f: c.sta, a: c.absolute, c: 4, s: 3},
		0x90: {f: c.bcc, a: c.relative, c: 1, s: 0},
		0x91: {f: c.sta, a: c.indirectY, c: 6, s: 2},
		0x94: {f: c.sty, a: c.zeropageX, c: 4, s: 2},
		0x95: {f: c.sta, a: c.zeropageX, c: 4, s: 2},
		0x98: {f: c.tya, a: c.implied, c: 2, s: 1},
		0x99: {f: c.sta, a: c.absoluteY, c: 5, s: 3},
		0x9a: {f: c.txs, a: c.implied, c: 2, s: 1},
		0x9d: {f: c.sta, a: c.absoluteX, c: 5, s: 3},
		0xa0: {f: c.ldy, a: c.immediate, c: 2, s: 2},
		0xa2: {f: c.ldx, a: c.immediate, c: 2, s: 2},
		0xa4: {f: c.ldy, a: c.zeropage, c: 3, s: 2},
		0xa5: {f: c.lda, a: c.zeropage, c: 3, s: 2},
		0xa6: {f: c.ldx, a: c.zeropage, c: 3, s: 2},
		0xa8: {f: c.tay, a: c.implied, c: 2, s: 1},
		0xa9: {f: c.lda, a: c.immediate, c: 2, s: 2},
		0xaa: {f: c.tax, a: c.implied, c: 2, s: 1},
		0xac: {f: c.ldy, a: c.absolute, c: 4, s: 3},
		0xad: {f: c.lda, a: c.absolute, c: 4, s: 3},
		0xae: {f: c.ldx, a: c.absolute, c: 4, s: 3},
		0xb0: {f: c.bcs, a: c.relative, c: 1, s: 0},
		0xb1: {f: c.lda, a: c.indirectY, c: 5, s: 2},
		0xb4: {f: c.ldy, a: c.zeropageX, c: 4, s: 2},
		0xb6: {f: c.ldx, a: c.zeropageY, c: 4, s: 2},
		0xb9: {f: c.lda, a: c.absoluteY, c: 4, s: 3},
		0xba: {f: c.tsx, a: c.implied, c: 2, s: 1},
		0xbc: {f: c.ldy, a: c.absoluteX, c: 4, s: 3},
		0xbe: {f: c.ldx, a: c.absoluteY, c: 4, s: 3},
		0xc0: {f: c.cpy, a: c.immediate, c: 2, s: 2},
		0xc1: {f: c.cmp, a: c.indirectX, c: 6, s: 2},
		0xc4: {f: c.cpy, a: c.zeropage, c: 3, s: 2},
		0xc5: {f: c.cmp, a: c.zeropage, c: 3, s: 2},
		0xc6: {f: c.dec, a: c.zeropage, c: 5, s: 2},
		0xc8: {f: c.iny, a: c.implied, c: 2, s: 1},
		0xc9: {f: c.cmp, a: c.immediate, c: 2, s: 2},
		0xca: {f: c.dex, a: c.implied, c: 2, s: 1},
		0xcc: {f: c.cpy, a: c.absolute, c: 4, s: 3},
		0xcd: {f: c.cmp, a: c.absolute, c: 4, s: 3},
		0xd0: {f: c.bne, a: c.relative, c: 1, s: 0},
		0xd1: {f: c.cmp, a: c.indirectY, c: 5, s: 2},
		0xd5: {f: c.cmp, a: c.zeropageX, c: 4, s: 2},
		0xd8: {f: c.cld, a: c.implied, c: 2, s: 1},
		0xd9: {f: c.cmp, a: c.absoluteY, c: 4, s: 3},
		0xdd: {f: c.cmp, a: c.absoluteX, c: 4, s: 3},
		0xe0: {f: c.cpx, a: c.immediate, c: 2, s: 2},
		0xe5: {f: c.sbc, a: c.zeropage, c: 3, s: 2},
		0xe6: {f: c.inc, a: c.zeropage, c: 5, s: 2},
		0xe8: {f: c.inx, a: c.implied, c: 2, s: 1},
		0xea: {f: c.nop, a: c.implied, c: 2, s: 1},
		0xf0: {f: c.beq, a: c.relative, c: 1, s: 0},
		0xf8: {f: c.sed, a: c.implied, c: 2, s: 1},
		0xf9: {f: c.sbc, a: c.absoluteY, c: 4, s: 3},
	}
}

// Run is the main function that processes through the PRG ROM.
func (c *CPU) Run() {
	fmt.Println("CPU spawned, getting initial PC...")
	c.pc = Address(c.readBytes(0xfffc))

	fmt.Printf("Beginning execution loop at $%04x\n", c.pc)
	var opcode uint8
	executed := 0
	for {
		executed++
		opcode = c.readMem(c.pc)
		fmt.Printf("%d 0x%04x A:0x%02x X:0x%02x Y:0x%02x SP:0x%02x OP:%02x\n", executed, c.pc, c.a, c.x, c.y, c.sp, opcode)
		c.executeNext()

		// VBLANK
		if !c.vblank && c.cycleCount >= 27507 {
			c.startVBlank()
		}

		if c.vblank && c.cycleCount >= 29780 {
			c.endVBlank()
			c.cycleCount = 0
		}
	}
}

func (c *CPU) readMem(address Address) uint8 {
	var controlBus chan uint16

	switch {
	case address >= 0x8000:
		controlBus = c.cartridgeControlBus
	case address >= 0x2000 && address < 0x4000:
		controlBus = c.ppuControlBus
	case address < 0x2000:
		return c.ram[address%0x0800]
	default:
		log.Panicf("attempt to read invalid memory location: 0x%x", address)
	}

	controlBus <- uint16(address)
	c.readWriteBus <- 0 // read
	return <-c.dataBus
}

func (c *CPU) readBytes(address Address) uint16 {
	lowbyte := uint16(c.readMem(address))
	highbyte := uint16(c.readMem(address + 1))
	return highbyte<<8 | lowbyte
}

func (c *CPU) writeMem(address Address, val uint8) error {
	var controlBus chan uint16

	switch {
	case address >= 0x8000:
		controlBus = c.cartridgeControlBus
	case address >= 0x4016 && address < 0x4018:
		return fmt.Errorf("Controller functions not implemented")
	case address == 0x4015:
		fmt.Printf("Write to pAPU address: $%04x\n", address)
		return nil
	case address == 0x4014:
		return fmt.Errorf("OAMDMA not implemented")
	case address >= 0x4000 && address < 0x4014:
		fmt.Printf("Write to pAPU address: $%04x\n", address)
		return nil
	case address >= 0x2000 && address < 0x4000:
		controlBus = c.ppuControlBus
	case address >= 0x0000 && address < 0x2000:
		// 0x0800 bytes mirrored four times
		c.ram[address%0x0800] = val
		return nil
	default:
		return fmt.Errorf("attempt to write to an invalid memory location: 0x%x", address)
	}

	controlBus <- uint16(address)
	c.readWriteBus <- 1 // write
	c.dataBus <- val

	return nil
}

func (c *CPU) executeNext() {

	// Read next opcode at the PC
	opcode := c.readMem(c.pc)
	inst := c.operations[opcode]

	if inst.f == nil {
		log.Panicf("unimplemented opcode: $%02x", opcode)
	}

	inst.f(inst.a)
	c.cycleCount += inst.c
	c.pc = Address(uint16(c.pc) + inst.s)
}

func (c *CPU) stackPush(value uint8) {
	stackAddress := Address(0x100 + uint16(c.sp))
	c.writeMem(stackAddress, value)
	c.sp--
}

func (c *CPU) stackPop() uint8 {
	c.sp++
	stackAddress := Address(0x100 + uint16(c.sp))
	return c.readMem(stackAddress)
}

func (c *CPU) setNegative(value uint8) {
	c.negative = (value&(1<<7) != 0)
}

func (c *CPU) setZero(value uint8) {
	c.zero = (value == 0)
}

func (c *CPU) startVBlank() {
	c.vblank = true
	c.vblankBus <- true
}

func (c *CPU) endVBlank() {
	c.vblank = false
	c.vblankBus <- false
}

func binToBcd(val uint8) int8 {
	ones := int8(val & 0x0f)
	tens := int8(val & 0xf0 / 0x10)
	return (tens * 10) + ones
}

func page(address Address) int {
	return int(math.Floor(float64(address) / 0x4000))
}
