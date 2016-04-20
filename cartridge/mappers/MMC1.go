package mapper

type shiftregister struct {
	val uint8
}

// MMC1 is a mapper ASIC used in Nintendo's SxROM and NES-EVENT
// Game Pak boards. Most common SxROM boards are assigned to iNES Mapper 001.
type MMC1 struct {
	PRG     []byte
	pages   []byte
	control uint8
	load    shiftregister
	chr0    uint8
	chr1    uint8
	prg0    uint8
}

func (r shiftregister) read() uint8 {
	val := r.val
	r.reset()
	return val
}

func (r shiftregister) write(val uint8) bool {
	flag := r.val & 0x01
	r.val = r.val >> 1
	r.val = r.val | val&0x01<<4

	return flag == 1
}

func (r shiftregister) reset() {
	r.val = 1 << 4
}

// Init initializes the ROM pages and control register
func (r MMC1) Init() error {
	r.pages = r.PRG[0:8192]
	r.pages = append(r.pages, r.PRG[len(r.PRG)-8192:len(r.PRG)]...)
	r.control = 0xc0

	r.load.reset()
	return nil
}

// Read returns the value from ROM stored in the address specified
func (r MMC1) Read(address uint16) (byte, error) {
	relativeAddress := address - 0x8000
	return r.pages[relativeAddress], nil
}

// Write fills the load register, and writes to the register specified by the
// specified address if the load register is full.
func (r MMC1) Write(address uint16, value byte) error {
	if value&1<<7 != 0 {
		r.load.reset()
		return nil
	}

	full := r.load.write(value)
	if full {
		if address < 0xa000 {
			r.control = r.load.read()
		}
		if address >= 0xa000 && address < 0xc000 {
			r.chr0 = r.load.read()
		}
		if address >= 0xc000 && address < 0xe000 {
			r.chr1 = r.load.read()
		}
		if address >= 0xe000 {
			r.prg0 = r.load.read()
		}
	}

	return nil
}
