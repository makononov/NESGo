package mapper

// A Mapper maps cartridge ROM into CPU ROM.
type Mapper interface {
	Init(prg []byte) error
	Read(address uint16) (byte, error)
	Write(address uint16, value byte) error
}
