package gb

import "fmt"

type bus struct {
	gb *Gameboy
}

func newBus(gb *Gameboy) *bus {
	return &bus{
		gb: gb,
	}
}

func (b *bus) read(addr uint16) uint8 {
	if addr < 0x8000 {
		return b.gb.cart.read(addr)
	} else if addr < 0xA000 {
		return b.gb.ram.read(addr)
	} else if addr < 0xE000 {
		return b.gb.ram.read(addr)
	} else if addr < 0xFE00 {
		fmt.Println("!! WARNING !! Nintendo says use of this area is prohibited")
		return b.gb.ram.read(addr) // Nintendo says use of this area is prohibited
	} else if addr < 0xFEA0 {
		return b.gb.ram.read(addr)
	} else if addr < 0xFF00 {
		fmt.Println("!! WARNING !! Nintendo says use of this area is prohibited")
		return b.gb.ram.read(addr) // Nintendo says use of this area is prohibited
	} else if addr < 0xFF80 {
		if addr >= 0xFF04 && addr <= 0xFF07 {
			return b.gb.timer.read(addr)
		} else if addr == 0xFF46 {
			return b.gb.dma.value
		} else {
			return b.gb.ram.read(addr)
		}
	} else if addr < 0xFFFF {
		return b.gb.ram.read(addr) // High RAM
	} else if addr == 0xFFFF {
		return b.gb.ram.read(addr) // IE
	} else {
		panic(fmt.Sprintf("read out of range, addr: 0x%04X", addr))
	}
}

func (b *bus) write(addr uint16, value uint8) {
	if addr < 0x8000 {
		fmt.Printf("!! WARNING !! Writing to ROM, addr: %04X, value: %02X\n", addr, value)
		b.gb.cart.write(addr, value)
	} else if addr < 0xA000 {
		b.gb.ram.write(addr, value)
	} else if addr < 0xE000 {
		b.gb.ram.write(addr, value)
	} else if addr < 0xFE00 {
		fmt.Println("!! WARNING !! Nintendo says use of this area is prohibited")
		b.gb.ram.write(addr, value) // Nintendo says use of this area is prohibited
	} else if addr < 0xFEA0 {
		b.gb.ram.write(addr, value)
	} else if addr < 0xFF00 {
		fmt.Println("!! WARNING !! Nintendo says use of this area is prohibited")
		b.gb.ram.write(addr, value) // Nintendo says use of this area is prohibited
	} else if addr < 0xFF80 {
		if addr >= 0xFF04 && addr <= 0xFF07 {
			b.gb.timer.write(addr, value)
		} else if addr == 0xFF46 {
			b.gb.dma.value = value
			b.gb.dma.start()
		} else {
			b.gb.ram.write(addr, value)
		}
	} else if addr < 0xFFFF {
		b.gb.ram.write(addr, value) // High RAM
	} else if addr == 0xFFFF {
		b.gb.ram.write(addr, value) // IE
	} else {
		panic(fmt.Sprintf("read out of range, addr: 0x%04X", addr))
	}
}

func (b *bus) read16(addr uint16) uint16 {
	dataLo := b.read(addr)
	dataHi := b.read(addr + 1)
	return uint16(dataHi)<<8 | uint16(dataLo)
}

func (b *bus) write16(addr uint16, value uint16) {
	lsb := uint8(value)
	msb := uint8(value >> 8)
	b.write(addr, lsb)
	b.write(addr+1, msb)
}
