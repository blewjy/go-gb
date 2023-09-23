package gb

import (
	"fmt"
)

type cart struct {
	gb *Gameboy

	rom []uint8
	ram []uint8

	title    string
	cgb      uint8
	cartType uint8
	romSize  uint8
	ramSize  uint8

	romBankSelect      uint64
	ramEnable          bool
	advancedBankSelect uint8
	advancedBanking    bool
}

func newCart(gb *Gameboy, rom []uint8) *cart {
	c := &cart{
		gb: gb,

		rom: rom,
		ram: make([]uint8, 131072),

		title:    string(rom[0x0134:0x0144]),
		cgb:      rom[0x0143],
		cartType: rom[0x0147],
		romSize:  rom[0x0148],
		ramSize:  rom[0x0149],

		romBankSelect:      1,
		ramEnable:          false,
		advancedBankSelect: 0,
		advancedBanking:    false,
	}

	fmt.Printf("title: %v\n", c.title)
	fmt.Printf("cgb: $%02x\n", c.cgb)
	fmt.Printf("cartType: $%02x\n", c.cartType)
	fmt.Printf("romSize: $%02x\n", c.romSize)
	fmt.Printf("ramSize: $%02x\n", c.ramSize)

	fmt.Printf("len(rom): %v\n", len(c.rom))

	return c
}

func (c *cart) updateRom(rom []byte) {
	c.rom = rom
}

func (c *cart) read(addr uint16) uint8 {
	switch c.cartType {
	case 0:
		return c.rom[addr]
	case 1:
		if addr < 0x4000 {
			if !c.advancedBanking {
				return c.rom[addr]
			}
			romBankNumber := uint64(c.advancedBankSelect) << 4 & c.romBankSelect
			finalAddr := 0x4000*(romBankNumber) + uint64(addr)
			return c.rom[finalAddr]

		} else if addr < 0x8000 {
			bitMask := uint64(1<<(c.romSize+1)) - 1
			offset := uint64(addr) - 0x4000
			finalAddr := 0x4000*(c.romBankSelect&bitMask) + offset
			return c.rom[finalAddr]
		} else if addr >= 0xA000 && addr < 0xC000 {
			if c.ramEnable {
				return c.ram[addr-0xA000]
			}
			return 0xFF
		} else {
			panic("cart read out of range")
		}
	default:
		panic("unsupported cart type")
	}

}

func (c *cart) write(addr uint16, value uint8) {
	if c.gb.testMode {
		c.rom[addr] = value
	}

	switch c.cartType {
	case 0x00:
	case 0x01:
		if addr < 0x2000 {
			if value&0xF == 0xA {
				c.ramEnable = true
			} else {
				c.ramEnable = false
			}
		} else if addr < 0x4000 {
			c.romBankSelect = (c.romBankSelect & 0xe0) | uint64(value&0x1F)
			if c.romBankSelect == 0x00 || c.romBankSelect == 0x20 || c.romBankSelect == 0x40 || c.romBankSelect == 0x60 {
				c.romBankSelect++
			}
		} else if addr < 0x6000 {
			c.advancedBankSelect = value & 0x3
		} else if addr < 0x8000 {
			if c.ramSize > 0x02 || c.romSize > 0x04 { // only advanced banking if ram > 8kb || rom > 512kb
				c.advancedBanking = value&0x1 > 0
			}
		} else if addr >= 0xA000 && addr < 0xC000 {
			if c.ramEnable {
				c.ram[addr-0xA000] = value
			}
		} else {
			panic("cart write out of range")
		}

	default:
		panic("unsupported cart type")
	}
}
