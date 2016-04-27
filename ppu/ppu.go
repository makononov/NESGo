package ppu

import (
	"errors"
	"fmt"
)

// PPU emulates the Picture Processing Unit of the NES
type PPU struct {
	controlBus   chan uint16
	readWriteBus chan int
	dataBus      chan uint8

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
}

// Init initializes a PPU struct with default values and the passed in
// bus channel.
func (p *PPU) Init(controlBus chan uint16, readWriteBus chan int, dataBus chan uint8) {
	p.controlBus = controlBus
	p.readWriteBus = readWriteBus
	p.dataBus = dataBus
}

func (p *PPU) writeMem(address uint16, value uint8) error {
	if address >= 0x2000 && address < 0x4000 { // register write channel
		// each register interface is mirrored every 8 bytes.
		registerNumber := (address - 0x2000) % 8
		switch registerNumber {
		case 0: // PPUCTRL
			p.writePPUCTRL(value)
			break
		case 1:
			p.writePPUMASK(value)
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
		address := <-p.controlBus
		if <-p.readWriteBus == 0 { // read
			panic(errors.New("PPU Reading not implemented"))
		} else { // write
			p.writeMem(address, <-p.dataBus)
		}

	}
}

func (p *PPU) writePPUCTRL(val uint8) {
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

func (p *PPU) writePPUMASK(val uint8) {
	p.grayscale = (val&0x01 != 0)
	p.showLeftBackground = (val&0x02 != 0)
	p.showLeftSprites = (val&0x04 != 0)
	p.showBackground = (val&0x08 != 0)
	p.showSprites = (val&0x10 != 0)
	p.emphasizeRed = (val&0x20 != 0)
	p.emphasizeGreen = (val&0x40 != 0)
	p.emphasizeRed = (val&0x80 != 0)
}
