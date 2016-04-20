package mapper

import (
	"errors"
)

// NROM is a simple ROM mapper with no logic controller.
type NROM struct {
	PRG   []byte
	pages []byte
}

// Init initializes the memory pages, duplicating an 8KB ROM in both banks or
// splitting a 16KB ROM into two banks.
func (r NROM) Init() error {
	if len(r.PRG) < 8192 {
		return errors.New("Attempted to initialize mapper with invalid ROM data")
	}

	r.pages = r.PRG[0:8192]
	r.pages = append(r.pages, r.PRG[len(r.PRG)-8192:len(r.PRG)]...)
	return nil
}

// Read returns the data stored at the specified address.
func (r NROM) Read(address uint16) (byte, error) {
	relativeAddress := address - 0x8000
	return r.pages[relativeAddress], nil
}

// Write is not supported on NROM, and should return an error.
func (r NROM) Write(address uint16, value byte) error {
	return errors.New("NROM does not support writing")
}
