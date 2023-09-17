package gb

type ram struct {
	memory []uint8
}

func newRam() *ram {
	return &ram{
		memory: make([]uint8, 65536),
	}
}

func (r *ram) read(addr uint16) uint8 {
	return r.memory[addr-0x8000]
}

func (r *ram) write(addr uint16, value uint8) {
	r.memory[addr-0x8000] = value
}
