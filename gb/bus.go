package gb

import "fmt"

type bus struct{}

func newBus() *bus {
	return &bus{}
}

func (b *bus) read(addr uint16) uint8 {
	if addr < 0x8000 {
		return gb.cart.read(addr)
	} else if addr < 0xA000 {
		return gb.ram.read(addr)
	} else if addr < 0xE000 {
		return gb.ram.read(addr)
	} else if addr < 0xFE00 {
		fmt.Println("!! WARNING !! Nintendo says use of this area is prohibited")
		return gb.ram.read(addr) // Nintendo says use of this area is prohibited
	} else if addr < 0xFEA0 {
		return gb.ram.read(addr)
	} else if addr < 0xFF00 {
		fmt.Println("!! WARNING !! Nintendo says use of this area is prohibited")
		return gb.ram.read(addr) // Nintendo says use of this area is prohibited
	} else if addr < 0xFF80 {
		return gb.ram.read(addr)
	} else if addr < 0xFFFF {
		return gb.ram.read(addr) // High RAM
	} else if addr == 0xFFFF {
		return gb.ram.read(addr) // IE
	} else {
		panic(fmt.Sprintf("read out of range, addr: 0x%04X", addr))
	}
}

func (b *bus) write(addr uint16, value uint8) {
	if addr < 0x8000 {
		gb.cart.write(addr, value)
	} else if addr < 0xA000 {
		gb.ram.write(addr, value)
	} else if addr < 0xE000 {
		gb.ram.write(addr, value)
	} else if addr < 0xFE00 {
		fmt.Println("!! WARNING !! Nintendo says use of this area is prohibited")
		gb.ram.write(addr, value) // Nintendo says use of this area is prohibited
	} else if addr < 0xFEA0 {
		gb.ram.write(addr, value)
	} else if addr < 0xFF00 {
		fmt.Println("!! WARNING !! Nintendo says use of this area is prohibited")
		gb.ram.write(addr, value) // Nintendo says use of this area is prohibited
	} else if addr < 0xFF80 {
		gb.ram.write(addr, value)
	} else if addr < 0xFFFF {
		gb.ram.write(addr, value) // High RAM
	} else if addr == 0xFFFF {
		gb.ram.write(addr, value) // IE
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
