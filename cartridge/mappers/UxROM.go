package mapper

import (
	"errors"
	"fmt"
)

// UxROM is a simple mapper with one siwtchable and one fixed page
type UxROM struct {
	PRG  []byte
	page int
}

// Init implements mapper.Init()
func (r *UxROM) Init(prg []byte) error {
	fmt.Println("Loaded mapper UxROM")
	r.PRG = prg
	r.page = 0

	return nil
}

// Read implements mapper.Read()
func (r *UxROM) Read(address uint16) (byte, error) {
	var offset uint16
	var pageStart int

	if address >= 0xc000 {
		offset = address - 0xc000
		pageStart = len(r.PRG) - 0x4000
	} else {
		offset = address - 0x8000
		pageStart = r.page * 0x4000
	}

	return r.PRG[pageStart+int(offset)], nil
}

// Write implements mapper.Write()
func (r *UxROM) Write(address uint16, value byte) error {
	if address < 0x8000 {
		return errors.New("Mapper write to invalid address")
	}

	r.page = int(value)
	return nil
}
