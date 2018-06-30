package ppu

import "fmt"

// PPU emulates the Picture Processing Unit of the NES
type PPU struct {
	controlBus   chan uint16
	readWriteBus chan int
	dataBus      chan uint8
	vblankBus    chan bool

	// PPUCTRL flags
	baseNametableAddress          uint16
	vramAddressIncrement          int
	spritePatternTableAddress     uint16
	backgroundPatternTableAddress uint16
	doubleHeightSprites           bool
	ppuMaster                     bool
	vblankNMI                     bool

	// PPUMASK flags
	grayscale          bool
	showLeftBackground bool
	showLeftSprites    bool
	showBackground     bool
	showSprites        bool
	emphasizeRed       bool
	emphasizeGreen     bool
	emphasizeBlue      bool

	// PPUSTATUS flags
	vblank         bool
	sprite0Hit     bool
	spriteOverflow bool
}

// Init initializes a PPU struct with default values and the passed in
// bus channel.
func (p *PPU) Init(controlBus chan uint16, readWriteBus chan int, dataBus chan uint8, vblankBus chan bool) {
	p.controlBus = controlBus
	p.readWriteBus = readWriteBus
	p.dataBus = dataBus
	p.vblankBus = vblankBus
}

func (p *PPU) readMem(address uint16) (uint8, error) {
	if address == 0x2002 {
		return p.ppuSTATUS(), nil
	}

	return 0, fmt.Errorf("Attempt to read PPU memory at 0x%x - not implemented", address)
}

func (p *PPU) writeMem(address uint16, value uint8) error {
	if address >= 0x2000 && address < 0x4000 { // register write channel
		// each register interface is mirrored every 8 bytes.
		registerNumber := (address - 0x2000) % 8
		switch registerNumber {
		case 0: // PPUCTRL
			p.setPPUCTRL(value)
			break
		case 1:
			p.setPPUMASK(value)
			break
		default:
			return fmt.Errorf("Attempt to write to PPU Register #%d - not implemented", registerNumber)
		}
		return nil
	}
	return fmt.Errorf("Attempt to write to PPU memory at 0x%x - not implemented", address)
}

// Run executes the threaded main loop of the PPU
func (p *PPU) Run() {
	fmt.Println("PPU spawned, awaiting instructions...")
	for {
		select {
		case p.vblank = <-p.vblankBus:
			if p.vblank {
				fmt.Println("*** VBLANK ***")
				// VBLANK CODE
			}

		case address := <-p.controlBus:
			if <-p.readWriteBus == 0 { // read
				val, err := p.readMem(address)
				if err != nil {
					panic(err)
				}
				p.dataBus <- val
			} else { // write
				if err := p.writeMem(address, <-p.dataBus); err != nil {
					panic(err)
				}
			}
		}
	}
}

func (p *PPU) setPPUCTRL(val uint8) {
	p.baseNametableAddress = 0x2000 + (uint16(val&0x03) * 0x0400)

	if val&0x04 != 0 {
		p.vramAddressIncrement = 32
	} else {
		p.vramAddressIncrement = 1
	}

	if val&0x08 != 0 {
		p.spritePatternTableAddress = 0x1000
	} else {
		p.spritePatternTableAddress = 0x0000
	}

	if val&0x10 != 0 {
		p.backgroundPatternTableAddress = 0x1000
	} else {
		p.backgroundPatternTableAddress = 0x0000
	}

	p.doubleHeightSprites = (val&0x20 != 0)

	p.ppuMaster = (val&0x40 == 0)
	p.vblankNMI = (val&0x80 != 0)
}

func (p *PPU) setPPUMASK(val uint8) {
	p.grayscale = (val&0x01 != 0)
	p.showLeftBackground = (val&0x02 != 0)
	p.showLeftSprites = (val&0x04 != 0)
	p.showBackground = (val&0x08 != 0)
	p.showSprites = (val&0x10 != 0)
	p.emphasizeRed = (val&0x20 != 0)
	p.emphasizeGreen = (val&0x40 != 0)
	p.emphasizeRed = (val&0x80 != 0)
}

func (p *PPU) ppuSTATUS() uint8 {
	value := uint8(0)
	if p.vblank {
		value |= 1 << 7
	}
	if p.sprite0Hit {
		value |= 1 << 6
	}
	if p.spriteOverflow {
		value |= 1 << 5
	}
	return value
}
