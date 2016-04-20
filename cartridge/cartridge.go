package cartridge

import (
	"fmt"

	"github.com/makononov/NESGo/cartridge/mappers"
)

const prgRomBlockSize int = 16384
const chrRomBlockSize int = 8192
const (
	// MirrorHorizontal is used to indicate that the cartridge uses
	// horizontal mirroring
	MirrorHorizontal = iota

	// MirrorVertical is used to indicate that the cartridge uses
	// vertical mirroring
	MirrorVertical = iota
)

const (
	// NTSC indicates the cartridge uses the NTCS TV system format
	NTSC = iota

	// PAL indicates the cartridge uses the PAL TV system format
	PAL = iota
)

// A Cartridge represents a game cartridge loaded into the system. It is
// generally created from a ROM file by calling cartridge.ParseROM().
type Cartridge struct {
	PrgRomSize        int
	ChrRomSize        int
	MapperID          int
	FourScreen        bool
	TrainerPresent    bool
	BatteryBackedSRAM bool
	Mirroring         int
	Playchoice10      bool
	VsUnisystem       bool
	TVSystemFormat    int

	Mapper mapper.Mapper

	Trainer []byte
	CHR     []byte
	RAM     []byte
}

// SetPrgRomSize sets the program ROM size of the cartridge, taking in to
// account the block size.
func (cartridge *Cartridge) SetPrgRomSize(size int) {
	cartridge.PrgRomSize = size * prgRomBlockSize
}

// SetChrRomSize sets the character ROM size of the cartridge, taking in
// to account the block size.
func (cartridge *Cartridge) SetChrRomSize(size int) {
	cartridge.ChrRomSize = size * chrRomBlockSize
}

// Read returns a byte located at the passed in address.
func (cartridge *Cartridge) Read(address uint16) (byte, error) {
	if address < 0x8000 {
		return byte(0), fmt.Errorf("ROM address out of range: %x", address)
	}

	return cartridge.Mapper.Read(address)
}
