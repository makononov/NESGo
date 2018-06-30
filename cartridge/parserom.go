package cartridge

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/makononov/NESGo/cartridge/mappers"
)

// ParseROM parses a ROM file and returns a cartridge object for use by the
// system.
func ParseROM(filename string) (*Cartridge, error) {
	fmt.Printf("Parsing %s...\n", filename)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cart := new(Cartridge)

	// Verify magic number
	if string(data[0:4]) != "NES\u001a" {
		return nil, fmt.Errorf("Invalid magic number in ROM file: %s", hex.Dump(data[0:4]))
	}

	cart.SetPrgRomSize(int(data[4]))
	fmt.Printf("Found ROM size: %d\n", cart.PrgRomSize)
	cart.SetChrRomSize(int(data[5]))
	fmt.Printf("Found CHR size: %d\n", cart.ChrRomSize)

	if err = parseFlags6(cart, data[6]); err != nil {
		return nil, err
	}

	var nes2 bool
	if nes2, err = parseFlags7(cart, data[7]); err != nil {
		return nil, err
	}

	if nes2 {
		return nil, errors.New("This in an NES2.0 ROM, which is currently not supported.")
	}

	if cart.Playchoice10 {
		return nil, errors.New("This is a Playchoice 10 ROM, which is currently not supported.")
	}

	if cart.VsUnisystem {
		return nil, errors.New("This is a Vs Unisystem ROM, which is currently not supported.")
	}

	// Verify bytes 11-15 are zeroed
	for _, num := range data[11:16] {
		if int(num) != 0 {
			fmt.Printf("ROM containes junk code in bytes 11-15: '%s'\n", string(data[10:16]))
			break
		}
	}

	if err = parseFlags9(cart, data[9]); err != nil {
		return nil, err
	}

	position := 16
	if cart.TrainerPresent {
		cart.Trainer = make([]byte, 512)
		copy(cart.Trainer, data[position:position+512])
		position = position + 512
	}

	// Initialize mapper
	// Copy the data so the file can close
	prg := make([]byte, cart.PrgRomSize)
	copy(prg, data[position:position+cart.PrgRomSize])

	switch cart.MapperID {
	case 0:
		cart.Mapper = new(mapper.NROM)
	case 1:
		cart.Mapper = new(mapper.MMC1)
	case 2:
		cart.Mapper = new(mapper.UxROM)
	default:
		return nil, fmt.Errorf("Mapper %d not yet implemented", cart.MapperID)
	}
	cart.Mapper.Init(prg)
	position = position + cart.PrgRomSize

	cart.CHR = make([]byte, cart.ChrRomSize)
	copy(cart.CHR, data[position:position+cart.ChrRomSize])
	position = position + cart.ChrRomSize

	return cart, nil
}

func parseFlags6(cart *Cartridge, flagByte byte) error {
	flag := int(flagByte)
	// bits 0-4 are the low nibble of the mapper ID
	cart.MapperID = (cart.MapperID & 0xf0) | (flag&0xf0)>>4
	cart.FourScreen = flag&0x08 == 1
	cart.TrainerPresent = flag&0x04 == 1
	cart.BatteryBackedSRAM = flag&0x02 == 1

	if flag&0x01 == 0 {
		cart.Mirroring = MirrorHorizontal
	} else {
		cart.Mirroring = MirrorVertical
	}

	return nil
}

func parseFlags7(cart *Cartridge, flagByte byte) (bool, error) {
	flag := int(flagByte)
	cart.MapperID = (cart.MapperID & 0x0f) | (flag & 0xf0)
	cart.Playchoice10 = flag&0x02 == 1
	cart.VsUnisystem = flag&0x01 == 1

	return flag&0x0c == 0x08, nil
}

func parseFlags9(cart *Cartridge, flagByte byte) error {
	flag := int(flagByte)
	if flag&0xfe != 0 {
		return errors.New("Invalid file format (flag 9)")
	}

	if flag&0x01 == 0 {
		cart.TVSystemFormat = NTSC
	} else {
		cart.TVSystemFormat = PAL
	}

	return nil
}
