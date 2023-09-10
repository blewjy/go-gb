package gb

import (
	"fmt"
)

type register uint8

const (
	reg_a register = iota
	reg_b
	reg_c
	reg_d
	reg_e
	reg_h
	reg_l
	reg_af
	reg_bc
	reg_de
	reg_hl
	reg_sp
)

func (r register) name() string {
	switch r {
	case reg_a:
		return "A"
	case reg_b:
		return "B"
	case reg_c:
		return "C"
	case reg_d:
		return "D"
	case reg_e:
		return "E"
	case reg_h:
		return "H"
	case reg_l:
		return "L"
	case reg_af:
		return "AF"
	case reg_bc:
		return "BC"
	case reg_de:
		return "DE"
	case reg_hl:
		return "HL"
	case reg_sp:
		return "SP"
	default:
		panic(fmt.Sprintf("no such register %v", r))
	}
}

func (c *cpu) readRegister(r register) uint8 {
	switch r {
	case reg_a:
		return c.a
	case reg_b:
		return c.b
	case reg_c:
		return c.c
	case reg_d:
		return c.d
	case reg_e:
		return c.e
	case reg_h:
		return c.h
	case reg_l:
		return c.l
	default:
		panic(fmt.Sprintf("no such register %v", r))
	}
}

func (c *cpu) readRegister16(r register) uint16 {
	switch r {
	case reg_af:
		return uint16(c.a)<<8 | uint16(c.f)
	case reg_bc:
		return uint16(c.b)<<8 | uint16(c.c)
	case reg_de:
		return uint16(c.d)<<8 | uint16(c.e)
	case reg_hl:
		return uint16(c.h)<<8 | uint16(c.l)
	case reg_sp:
		return c.sp
	default:
		panic(fmt.Sprintf("no such register %v", r))
	}
}

func (c *cpu) writeRegister(r register, v uint8) {
	switch r {
	case reg_a:
		c.a = v
	case reg_b:
		c.b = v
	case reg_c:
		c.c = v
	case reg_d:
		c.d = v
	case reg_e:
		c.e = v
	case reg_h:
		c.h = v
	case reg_l:
		c.l = v
	default:
		panic(fmt.Sprintf("no such register %v", r))
	}
}

func (c *cpu) writeRegister16(r register, v uint16) {
	switch r {
	case reg_af:
		c.a = uint8(v >> 8)
		c.f = uint8(v)
	case reg_bc:
		c.b = uint8(v >> 8)
		c.c = uint8(v)
	case reg_de:
		c.d = uint8(v >> 8)
		c.e = uint8(v)
	case reg_hl:
		c.h = uint8(v >> 8)
		c.l = uint8(v)
	case reg_sp:
		c.sp = v
	default:
		panic(fmt.Sprintf("no such register %v", r))
	}
}

func (c *cpu) push(v uint8) {
	c.sp -= 1
	gb.bus.write(c.sp, v)
}

func (c *cpu) pop() uint8 {
	v := gb.bus.read(c.sp)
	c.sp += 1
	return v
}

func (c *cpu) push16(v uint16) {
	c.sp -= 1
	gb.bus.write(c.sp, uint8(v>>8))
	c.sp -= 1
	gb.bus.write(c.sp, uint8(v))
}

func (c *cpu) pop16() uint16 {
	lo := gb.bus.read(c.sp)
	c.sp += 1
	hi := gb.bus.read(c.sp)
	c.sp += 1
	return uint16(hi)<<8 | uint16(lo)
}
