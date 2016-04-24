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

	controlRegister uint8
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
			p.controlRegister = value
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
