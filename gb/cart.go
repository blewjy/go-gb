package gb

type cart struct {
	rom []uint8
}

func newCart(rom []uint8) *cart {
	return &cart{
		rom: rom,
	}
}

func (c *cart) updateRom(rom []byte) {
	c.rom = rom
}

func (c *cart) read(addr uint16) uint8 {
	return c.rom[addr]
}

func (c *cart) write(addr uint16, value uint8) {
	c.rom[addr] = value
}
