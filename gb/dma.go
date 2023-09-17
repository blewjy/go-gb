package gb

type dma struct {
	gb *Gameboy

	value  uint8
	offset uint8
	active bool
}

func newDma(gb *Gameboy) *dma {
	return &dma{
		gb: gb,
	}
}

func (d *dma) start() {
	d.active = true
}

func (d *dma) tick() {
	if !d.active {
		return
	}

	v := d.gb.bus.read(uint16(d.value)<<8 + uint16(d.offset))
	d.gb.ram.write(0xFE00+uint16(d.offset), v)
	d.offset += 1
	if d.offset >= 0xA0 {
		d.active = false
		d.offset = 0
	}
}
