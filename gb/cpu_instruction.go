package gb

func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func (c *cpu) nop() uint8 {
	return 4
}

func (c *cpu) illegal() uint8 {
	panic("illegal instruction")
}

func (c *cpu) cb_prefix() uint8 {
	cbOpcode := gb.bus.read(c.pc)
	c.pc += 1
	return opcodeToCBInst[cbOpcode](c)
}

/** 8-bit load instructions **/

// ld_r_r: Load to the 8-bit register dr, data from the 8-bit register sr.
//
// Clock cycles: 4
// Opcodes:
func (c *cpu) ld_r_r(dr, sr register) uint8 {
	v := c.readRegister(sr)
	c.writeRegister(dr, v)
	// log.Info("ld_r_r\tLD %s, %s[%02X]", dr.name(), sr.name(), v)
	return 4
}

// ld_r_n: Load to the 8-bit register dr, the immediate data n.
//
// Clock cycles: 8
// Opcodes:
func (c *cpu) ld_r_n(dr register) uint8 {
	v := gb.bus.read(c.pc)
	c.pc += 1
	c.writeRegister(dr, v)
	// log.Info("ld_r_n\tLD %s, %02x", dr.name(), v)
	return 8
}

// ld_r_HL: Load to the 8-bit register dr, data from the absolute address specified by the 16-bit register HL.
//
// Clock cycles: 8
// Opcodes:
func (c *cpu) ld_r_HL(dr register) uint8 {
	a := c.readRegister16(reg_hl)
	v := gb.bus.read(a)
	c.writeRegister(dr, v)
	// log.Info("ld_r_HL\tLD %s, (HL)[%02X]", dr.name(), v)
	return 8
}

// ld_HL_r: Load to the absolute address specified by the 16-bit register HL, data from the 8-bit register sr.
//
// Clock cycles: 8
// Opcodes:
func (c *cpu) ld_HL_r(sr register) uint8 {
	v := c.readRegister(sr)
	a := c.readRegister16(reg_hl)
	gb.bus.write(a, v)
	// log.Info("ld_HL_r\tLD (HL), %s[%02X]", sr.name(), v)
	return 8
}

// ld_HL_n: Load to the absolute address specified by the 16-bit register HL, the immediate data n.
//
// Clock cycles: 12
// Opcodes:
func (c *cpu) ld_HL_n() uint8 {
	v := gb.bus.read(c.pc)
	c.pc += 1
	a := c.readRegister16(reg_hl)
	gb.bus.write(a, v)
	// log.Info("ld_HL_n\tLD (HL), %02X", v)
	return 12
}

// ld_a_BC: Load to the 8-bit A register, data from the absolute address specified by the 16-bit register BC.
//
// Clock cycles: 8
// Opcodes:
func (c *cpu) ld_a_BC() uint8 {
	a := c.readRegister16(reg_bc)
	v := gb.bus.read(a)
	c.writeRegister(reg_a, v)
	// log.Info("ld_a_BC\tLD A, (BC)[%02X]", v)
	return 8
}

// ld_a_DE: Load to the 8-bit A register, data from the absolute address specified by the 16-bit register DE.
//
// Clock cycles: 8
func (c *cpu) ld_a_DE() uint8 {
	a := c.readRegister16(reg_de)
	v := gb.bus.read(a)
	c.writeRegister(reg_a, v)
	// log.Info("ld_a_DE\tLD A, (DE)[%02X]", v)
	return 8
}

// ld_a_NN: Load to the 8-bit A register, data from the absolute address specified by the 16-bit operand nn.
//
// Clock cycles: 16
func (c *cpu) ld_a_NN() uint8 {
	a := gb.bus.read16(c.pc)
	c.pc += 2
	v := gb.bus.read(a)
	c.writeRegister(reg_a, v)
	// log.Info("ld_a_NN\tLD A, (%04X)[%02X]", a, v)
	return 16
}

// ld_BC_a: Load to the absolute address specified by the 16-bit register BC, data from the 8-bit A register.
//
// Clock cycles: 8
func (c *cpu) ld_BC_a() uint8 {
	a := c.readRegister16(reg_bc)
	v := c.readRegister(reg_a)
	gb.bus.write(a, v)
	// log.Info("ld_BC_a\tLD (%04X), A[%02X]", a, v)
	return 8
}

// ld_DE_a: Load to the absolute address specified by the 16-bit register DE, data from the 8-bit A register.
//
// Clock cycles: 8
func (c *cpu) ld_DE_a() uint8 {
	a := c.readRegister16(reg_de)
	v := c.readRegister(reg_a)
	gb.bus.write(a, v)
	// log.Info("ld_DE_a\tLD (%04X), A[%02X]", a, v)
	return 8
}

// ld_NN_a: Load to the absolute address specified by the 16-bit operand nn, data from the 8-bit A register.
//
// Clock cycles: 16
func (c *cpu) ld_NN_a() uint8 {
	a := gb.bus.read16(c.pc)
	c.pc += 2
	v := c.readRegister(reg_a)
	gb.bus.write(a, v)
	// log.Info("ld_NN_a\tLD (%04X), A[%02X]", a, v)
	return 16
}

// ld_a_ff00n: Load to the 8-bit A register, data from the address specified by the 8-bit immediate data n.
// The full 16-bit absolute address is obtained by setting the most significant byte to 0xFF and the least
// significant byte to the value of n, so the possible range is 0xFF00-0xFFFF.
//
// Clock cycles: 12
func (c *cpu) ld_a_ff00n() uint8 {
	n := gb.bus.read(c.pc)
	c.pc += 1
	a := 0xFF00 + uint16(n)
	v := gb.bus.read(a)
	c.writeRegister(reg_a, v)
	// log.Info("ld_a_ff00n\tLD A, (FF00+%02X)[%02X]", n, v)
	return 12
}

// ld_ff00n_a: Load to the address specified by the 8-bit immediate data n, data from the 8-bit A register. The full 16-bit
// absolute address is obtained by setting the most significant byte to 0xFF and the least significant byte to the
// value of n, so the possible range is 0xFF00-0xFFFF.
//
// Clock cycles: 12
func (c *cpu) ld_ff00n_a() uint8 {
	n := gb.bus.read(c.pc)
	c.pc += 1
	a := 0xFF00 + uint16(n)
	v := c.readRegister(reg_a)
	gb.bus.write(a, v)
	// log.Info("ld_ff00n_a\tLD (FF00+%02X), A[%02X]", n, v)
	return 12
}

// ld_a_ff00c: Load to the 8-bit A register, data from the address specified by the 8-bit C register. The full 16-bit absolute
// address is obtained by setting the most significant byte to 0xFF and the least significant byte to the value of C,
// so the possible range is 0xFF00-0xFFFF.
//
// Clock cycles: 8
func (c *cpu) ld_a_ff00c() uint8 {
	a := 0xFF00 + uint16(c.readRegister(reg_c))
	v := gb.bus.read(a)
	c.writeRegister(reg_a, v)
	// log.Info("ld_a_ff00c\tLD A, (FF00+C)[%02X]", v)
	return 8
}

// ld_ff00c_a: Load to the address specified by the 8-bit C register, data from the 8-bit A register. The full 16-bit absolute
// address is obtained by setting the most significant byte to 0xFF and the least significant byte to the value of C,
// so the possible range is 0xFF00-0xFFFF.
//
// Clock cycles: 8
func (c *cpu) ld_ff00c_a() uint8 {
	a := 0xFF00 + uint16(c.readRegister(reg_c))
	v := c.readRegister(reg_a)
	gb.bus.write(a, v)
	// log.Info("ld_ff00c_a\tLD (FF00+C), A[%02X]", v)
	return 8
}

// ldi_HL_a: Load to the absolute address specified by the 16-bit register HL, data from the 8-bit A register. The value of
// HL is incremented after the memory write.
//
// Clock cycles: 8
func (c *cpu) ldi_HL_a() uint8 {
	a := c.readRegister16(reg_hl)
	v := c.readRegister(reg_a)
	gb.bus.write(a, v)
	c.writeRegister16(reg_hl, a+1)
	// log.Info("ldi_hl_a\tLDI (HL), A[%02X]", v)
	return 8
}

// ldi_a_HL: Load to the 8-bit A register, data from the absolute address specified by the 16-bit register HL. The value of
// HL is incremented after the memory read.
//
// Clock cycles: 8
func (c *cpu) ldi_a_HL() uint8 {
	a := c.readRegister16(reg_hl)
	v := gb.bus.read(a)
	c.writeRegister(reg_a, v)
	c.writeRegister16(reg_hl, a+1)
	// log.Info("ldi_a_hl\tLDI A, (HL)[%02X]", v)
	return 8
}

// ldd_HL_a: Load to the absolute address specified by the 16-bit register HL, data from the 8-bit A register. The value of
// HL is decremented after the memory write.
//
// Clock cycles: 8
func (c *cpu) ldd_HL_a() uint8 {
	a := c.readRegister16(reg_hl)
	v := c.readRegister(reg_a)
	gb.bus.write(a, v)
	c.writeRegister16(reg_hl, a-1)
	// log.Info("ldd_HL_a\tLDD (HL), A[%02X]", v)
	return 8
}

// ldd_a_HL: Load to the 8-bit A register, data from the absolute address specified by the 16-bit register HL. The value of
// HL is decremented after the memory read.
//
// Clock cycles: 8
func (c *cpu) ldd_a_HL() uint8 {
	a := c.readRegister16(reg_hl)
	v := gb.bus.read(a)
	c.writeRegister(reg_a, v)
	c.writeRegister16(reg_hl, a-1)
	// log.Info("ldd_a_HL\tLDD A, (HL)[%02X]", v)
	return 8
}

/** 16-bit load instructions **/

// ld_rr_nn: Load to the 16-bit register rr, the immediate 16-bit data nn.
//
// Clock cycles: 12
func (c *cpu) ld_rr_nn(dr register) uint8 {
	v := gb.bus.read16(c.pc)
	c.pc += 2
	c.writeRegister16(dr, v)
	// log.Info("ld_rr_nn\tLD %s, %04X", dr.name(), v)
	return 12
}

// ld_NN_sp: Load to the absolute address specified by the 16-bit operand nn, data from the 16-bit SP register.
//
// Clock cycles: 20
func (c *cpu) ld_NN_sp() uint8 {
	v := c.sp
	a := gb.bus.read16(c.pc)
	c.pc += 2
	gb.bus.write16(a, v)
	// log.Info("ld_NN_sp\tLD (%04X), SP[%04X]", a, v)
	return 20
}

// ld_sp_hl: Load to the 16-bit SP register, data from the 16-bit HL register.
//
// Clock cycles: 8
func (c *cpu) ld_sp_hl() uint8 {
	v := c.readRegister16(reg_hl)
	c.sp = v
	// log.Info("ld_sp_hl\tLD SP, HL[%04X]", v)
	return 8
}

// push_rr: Push to the stack memory, data from the 16-bit register rr.
//
// Clock cycles: 16
func (c *cpu) push_rr(sr register) uint8 {
	v := c.readRegister16(sr)
	c.push16(v)
	// log.Info("push_rr\tPUSH %s[%04X]", sr.name(), v)
	return 16
}

// pop_rr: Pops to the 16-bit register rr, data from the stack memory.
// This instruction does not do calculations that affect flags, but POP AF completely replaces the F register
// value, so all flags are changed based on the 8-bit data that is read from memory.
//
// Clock cycles: 12
func (c *cpu) pop_rr(dr register) uint8 {
	v := c.pop16()
	c.writeRegister16(dr, v)

	// for POP AF, we need to set the lower 4 bits to 0
	c.f &= 0xF0

	return 12

	// log.Info("pop_rr\tPOP %s[%04X]", dr.name(), v)
}

/** 8-bit Arithmetic/Logic instructions **/

// add_a_r
//
// Clock cycles: 4
func (c *cpu) add_a_r(r register) uint8 {
	a := uint16(c.readRegister(reg_a))
	b := uint16(c.readRegister(r))
	result := a + b

	c.setFlag(flagZ, uint8(result) == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, (a&0x0F)+(b&0x0F) > 0x0F)
	c.setFlag(flagC, result > 0xFF)

	c.writeRegister(reg_a, uint8(result))

	// log.Info("add_a_r\tADD A, %s[%02X]", r.name(), b)
	return 4
}

// add_a_n
//
// Clock cycles: 8
func (c *cpu) add_a_n() uint8 {
	a := uint16(c.readRegister(reg_a))
	n := uint16(gb.bus.read(c.pc))
	c.pc += 1
	result := a + n

	c.setFlag(flagZ, uint8(result) == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, (a&0x0F)+(n&0x0F) > 0x0F)
	c.setFlag(flagC, result > 0xFF)

	c.writeRegister(reg_a, uint8(result))

	// log.Info("add_a_n\tADD A, 0x%02X", n)
	return 8
}

// add_a_HL
//
// Clock cycles: 8
func (c *cpu) add_a_HL() uint8 {
	a := uint16(c.readRegister(reg_a))
	hlValue := c.readRegister16(reg_hl)
	n := uint16(gb.bus.read(hlValue))
	result := a + n

	c.setFlag(flagZ, uint8(result) == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, (a&0x0F)+(n&0x0F) > 0x0F)
	c.setFlag(flagC, result > 0xFF)

	c.writeRegister(reg_a, uint8(result))

	return 8

	// log.Info("add_a_hl\tADD A, (HL)[%02X]", n)
}

// adc_a_r
//
// Clock cycles: 4
func (c *cpu) adc_a_r(r register) uint8 {
	a := uint16(c.readRegister(reg_a))
	b := uint16(c.readRegister(r))
	carry := c.getFlag(flagC)
	carryInt := uint16(0)
	if carry {
		carryInt = 1
	}
	result := a + b + carryInt

	c.setFlag(flagZ, uint8(result) == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, (a&0x0F)+(b&0x0F)+carryInt > 0x0F)
	c.setFlag(flagC, result > 0xFF)

	c.writeRegister(reg_a, uint8(result))

	return 4

	// log.Info("adc_a_r\tADC A, %s[%02X]", r.name(), b)
}

// adc_a_n
//
// Clock cycles: 8
func (c *cpu) adc_a_n() uint8 {
	a := uint16(c.readRegister(reg_a))
	n := uint16(gb.bus.read(c.pc))
	c.pc += 1

	carry := c.getFlag(flagC)
	carryInt := uint16(0)
	if carry {
		carryInt = 1
	}
	result := a + n + carryInt

	c.setFlag(flagZ, uint8(result) == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, (a&0x0F)+(n&0x0F)+carryInt > 0x0F)
	c.setFlag(flagC, result > 0xFF)

	c.writeRegister(reg_a, uint8(result))

	return 8

	// log.Info("adc_a_n\tADC A, 0x%02X", n)
}

// adc_a_HL
//
// Clock cycles: 8
func (c *cpu) adc_a_HL() uint8 {
	a := uint16(c.readRegister(reg_a))
	hlValue := c.readRegister16(reg_hl)
	n := uint16(gb.bus.read(hlValue))
	carry := c.getFlag(flagC)
	carryInt := uint16(0)
	if carry {
		carryInt = 1
	}
	result := a + n + carryInt

	c.setFlag(flagZ, uint8(result) == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, (a&0x0F)+(n&0x0F)+carryInt > 0x0F)
	c.setFlag(flagC, result > 0xFF)

	c.writeRegister(reg_a, uint8(result))

	return 8

	// log.Info("adc_a_hl\tADC A, (HL)[%02X]", n)
}

// sub_r
//
// Clock cycles: 4
func (c *cpu) sub_r(r register) uint8 {
	a := c.readRegister(reg_a)
	b := c.readRegister(r)
	result := a - b

	c.setFlag(flagZ, uint8(result) == 0)
	c.setFlag(flagN, true)
	c.setFlag(flagH, (a&0x0F) < (b&0x0F))
	c.setFlag(flagC, a < b)

	c.writeRegister(reg_a, result)

	return 4

	// log.Info("sub_r\tSUB A, %s[%02X]", r.name(), b)
}

// sub_n
//
// Clock cycles: 8
func (c *cpu) sub_n() uint8 {
	a := c.readRegister(reg_a)
	n := gb.bus.read(c.pc)
	c.pc += 1
	result := a - n

	c.setFlag(flagZ, uint8(result) == 0)
	c.setFlag(flagN, true)
	c.setFlag(flagH, (a&0x0F) < (n&0x0F))
	c.setFlag(flagC, a < n)

	c.writeRegister(reg_a, result)

	return 8

	// log.Info("sub_n\tSUB A, 0x%02X", n)
}

// sub_HL
//
// Clock cycles: 8
func (c *cpu) sub_HL() uint8 {
	a := c.readRegister(reg_a)
	hlValue := c.readRegister16(reg_hl)
	n := gb.bus.read(hlValue)
	result := a - n

	c.setFlag(flagZ, uint8(result) == 0)
	c.setFlag(flagN, true)
	c.setFlag(flagH, (a&0x0F) < (n&0x0F))
	c.setFlag(flagC, a < n)

	c.writeRegister(reg_a, result)

	return 8

	// log.Info("sub_hl\tSUB A, (HL)[%02X]", n)
}

// sbc_a_r
//
// Clock cycles: 4
func (c *cpu) sbc_a_r(r register) uint8 {
	a := c.readRegister(reg_a)
	b := c.readRegister(r)
	carry := c.getFlag(flagC)
	carryInt := uint8(0)
	if carry {
		carryInt = 1
	}
	result := a - b - carryInt

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, true)
	c.setFlag(flagH, (a&0x0F) < ((b&0x0F)+carryInt))
	c.setFlag(flagC, uint16(a) < uint16(b)+uint16(carryInt))

	c.writeRegister(reg_a, result)

	return 4

	// log.Info("sbc_a_r\tSBC A, %s[%02X]", r.name(), b)
}

// sbc_a_n
//
// Clock cycles: 8
func (c *cpu) sbc_a_n() uint8 {
	a := c.readRegister(reg_a)
	n := gb.bus.read(c.pc)
	c.pc += 1
	carry := c.getFlag(flagC)
	carryInt := uint8(0)
	if carry {
		carryInt = 1
	}

	result := a - n - carryInt

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, true)
	c.setFlag(flagH, (a&0x0F) < ((n&0x0F)+carryInt))
	c.setFlag(flagC, uint16(a) < uint16(n)+uint16(carryInt))

	c.writeRegister(reg_a, result)

	return 8

	// log.Info("sbc_a_n\tSBC A, 0x%02X", n)
}

// sbc_a_HL
//
// Clock cycles: 8
func (c *cpu) sbc_a_HL() uint8 {
	a := c.readRegister(reg_a)
	hlValue := c.readRegister16(reg_hl)
	n := gb.bus.read(hlValue)
	carry := c.getFlag(flagC)
	carryInt := uint8(0)
	if carry {
		carryInt = 1
	}
	result := a - n - carryInt

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, true)
	c.setFlag(flagH, (a&0x0F) < ((n&0x0F)+carryInt))
	c.setFlag(flagC, uint16(a) < uint16(n)+uint16(carryInt))

	c.writeRegister(reg_a, result)

	return 8

	// log.Info("sbc_a_hl\tSBC A, (HL)[%02X]", n)
}

// and_r
//
// Clock cycles: 4
func (c *cpu) and_r(r register) uint8 {
	a := c.readRegister(reg_a)
	b := c.readRegister(r)
	result := a & b

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, true)
	c.setFlag(flagC, false)

	c.writeRegister(reg_a, result)

	return 4

	// log.Info("and_r\tAND A, %s[%02X]", r.name(), b)
}

// and_n
//
// Clock cycles: 8
func (c *cpu) and_n() uint8 {
	a := c.readRegister(reg_a)
	n := gb.bus.read(c.pc)
	c.pc += 1
	result := a & n

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, true)
	c.setFlag(flagC, false)

	c.writeRegister(reg_a, result)

	return 8

	// log.Info("and_n\tAND A, 0x%02X", n)
}

// and_HL
//
// Clock cycles: 8
func (c *cpu) and_HL() uint8 {
	a := c.readRegister(reg_a)
	hlValue := c.readRegister16(reg_hl)
	n := gb.bus.read(hlValue)
	result := a & n

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, true)
	c.setFlag(flagC, false)

	c.writeRegister(reg_a, result)

	return 8

	// log.Info("and_hl\tAND A, (HL)[%02X]", n)
}

// xor_r
//
// Clock cycles: 4
func (c *cpu) xor_r(r register) uint8 {
	a := c.readRegister(reg_a)
	b := c.readRegister(r)
	result := a ^ b

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, false)

	c.writeRegister(reg_a, result)

	return 4

	// log.Info("xor_r\tXOR A, %s[%02X]", r.name(), b)
}

// xor_n
//
// Clock cycles: 8
func (c *cpu) xor_n() uint8 {
	a := c.readRegister(reg_a)
	n := gb.bus.read(c.pc)
	c.pc += 1
	result := a ^ n

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, false)

	c.writeRegister(reg_a, result)

	return 8

	// log.Info("xor_n\tXOR A, 0x%02X", n)
}

// xor_HL
//
// Clock cycles: 8
func (c *cpu) xor_HL() uint8 {
	a := c.readRegister(reg_a)
	hlValue := c.readRegister16(reg_hl)
	n := gb.bus.read(hlValue)
	result := a ^ n

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, false)

	c.writeRegister(reg_a, result)

	return 8

	// log.Info("xor_hl\tXOR A, (HL)[%02X]", n)
}

// or_r
//
// Clock cycles: 4
func (c *cpu) or_r(r register) uint8 {
	a := c.readRegister(reg_a)
	b := c.readRegister(r)
	result := a | b

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, false)

	c.writeRegister(reg_a, result)

	return 4

	// log.Info("or_r\tOR A, %s[%02X]", r.name(), b)
}

// or_n
//
// Clock cycles: 8
func (c *cpu) or_n() uint8 {
	a := c.readRegister(reg_a)
	n := gb.bus.read(c.pc)
	c.pc += 1
	result := a | n

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, false)

	c.writeRegister(reg_a, result)

	return 8

	// log.Info("or_n\tOR A, 0x%02X", n)
}

// or_HL
//
// Clock cycles: 8
func (c *cpu) or_HL() uint8 {
	a := c.readRegister(reg_a)
	hlValue := c.readRegister16(reg_hl)
	n := gb.bus.read(hlValue)
	result := a | n

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, false)

	c.writeRegister(reg_a, result)

	return 8

	// log.Info("or_hl\tOR A, (HL)[%02X]", n)
}

// cp_r
//
// Clock cycles: 4
func (c *cpu) cp_r(r register) uint8 {
	a := c.readRegister(reg_a)
	b := c.readRegister(r)

	c.setFlag(flagZ, a == b)
	c.setFlag(flagN, true)
	c.setFlag(flagH, (a&0x0F) < (b&0x0F))
	c.setFlag(flagC, a < b)

	return 4

	// log.Info("cp_r\tCP A, %s[%02X]", r.name(), b)
}

// cp_n
//
// Clock cycles: 8
func (c *cpu) cp_n() uint8 {
	a := c.readRegister(reg_a)
	n := gb.bus.read(c.pc)
	c.pc += 1

	c.setFlag(flagZ, a == n)
	c.setFlag(flagN, true)
	c.setFlag(flagH, (a&0x0F) < (n&0x0F))
	c.setFlag(flagC, a < n)

	return 8

	// log.Info("cp_n\tCP A, 0x%02X", n)
}

// cp_HL
//
// Clock cycles: 8
func (c *cpu) cp_HL() uint8 {
	a := c.readRegister(reg_a)
	hlValue := c.readRegister16(reg_hl)
	n := gb.bus.read(hlValue)

	c.setFlag(flagZ, a == n)
	c.setFlag(flagN, true)
	c.setFlag(flagH, (a&0x0F) < (n&0x0F))
	c.setFlag(flagC, a < n)

	return 8

	// log.Info("cp_hl\tCP A, (HL)[%02X]", n)
}

// inc_r
//
// Clock cycles: 4
func (c *cpu) inc_r(r register) uint8 {
	value := c.readRegister(r)
	result := value + 1

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, (value&0x0F) == 0x0F)

	c.writeRegister(r, result)

	return 4

	// log.Info("inc_r\tINC %s[%02X]", r.name(), result)
}

// inc_HL
//
// Clock cycles: 12
func (c *cpu) inc_HL() uint8 {
	hlValue := c.readRegister16(reg_hl)
	value := gb.bus.read(hlValue)
	result := value + 1

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, (value&0x0F) == 0x0F)

	gb.bus.write(hlValue, result)

	return 12

	// log.Info("inc_hl\tINC (HL)[%02X]", result)
}

// dec_r
//
// Clock cycles: 4
func (c *cpu) dec_r(r register) uint8 {
	value := c.readRegister(r)
	result := value - 1

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, true)
	c.setFlag(flagH, (value&0x0F) == 0x00)

	c.writeRegister(r, result)

	return 4

	// log.Info("dec_r\tDEC %s[%02X]", r.name(), result)
}

// dec_HL
//
// Clock cycles: 12
func (c *cpu) dec_HL() uint8 {
	hlValue := c.readRegister16(reg_hl)
	value := gb.bus.read(hlValue)
	result := value - 1

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, true)
	c.setFlag(flagH, (value&0x0F) == 0x00)

	gb.bus.write(hlValue, result)

	return 12

	// log.Info("dec_hl\tDEC (HL)[%02X]", result)
}

// daa
//
// Clock cycles: 4
func (c *cpu) daa() uint8 {
	a := c.readRegister(reg_a)
	u := uint8(0)
	fc := false
	if c.getFlag(flagH) || (!c.getFlag(flagN) && (a&0xF) > 9) {
		u = 6
	}

	if c.getFlag(flagC) || (!c.getFlag(flagN) && a > 0x99) {
		u |= 0x60
		fc = true
	}

	if c.getFlag(flagN) {
		a -= u
	} else {
		a += u
	}

	c.setFlag(flagZ, a == 0)
	c.setFlag(flagH, false)
	c.setFlag(flagC, fc)

	c.writeRegister(reg_a, a)

	return 4

	// log.Info("daa\tDAA A[%02X]", a)
}

// cpl
//
// Clock cycles: 4
func (c *cpu) cpl() uint8 {
	a := c.readRegister(reg_a)
	result := ^a

	c.setFlag(flagN, true)
	c.setFlag(flagH, true)

	c.writeRegister(reg_a, result)

	return 4

	// log.Info("cpl\tCPL A[%02X]", result)
}

/** 16-bit Arithmetic/Logic instructions **/

// add_hl_rr
//
// Clock cycles: 8
func (c *cpu) add_hl_rr(rr register) uint8 {
	hlValue := uint32(c.readRegister16(reg_hl))
	rrValue := uint32(c.readRegister16(rr))

	result := hlValue + rrValue

	c.setFlag(flagN, false)
	c.setFlag(flagH, ((hlValue&0x0FFF)+(rrValue&0x0FFF)) > 0x0FFF)
	c.setFlag(flagC, result > 0xFFFF)

	c.writeRegister16(reg_hl, uint16(result))

	return 8

	// log.Info("add_hl_rr\tADD HL, %s[%04X]", rr.name(), rrValue)
}

// inc_rr
//
// Clock cycles: 8
func (c *cpu) inc_rr(rr register) uint8 {
	value := c.readRegister16(rr)
	result := value + 1

	c.writeRegister16(rr, result)
	return 8
	// log.Info("inc_rr\tINC %s[%04X]", rr.name(), result)
}

// dec_rr
//
// Clock cycles: 8
func (c *cpu) dec_rr(rr register) uint8 {
	value := c.readRegister16(rr)
	result := value - 1

	c.writeRegister16(rr, result)
	return 8
	// log.Info("dec_rr\tDEC %s[%04X]", rr.name(), result)
}

// add_sp_dd
//
// Clock cycles: 16
func (c *cpu) add_sp_dd() uint8 {
	n := int8(gb.bus.read(c.pc))
	c.pc += 1

	oldSP := c.sp
	result := int32(oldSP) + int32(n)

	c.setFlag(flagZ, false)
	c.setFlag(flagN, false)
	c.setFlag(flagH, ((oldSP&0x0F)+(uint16(n)&0x0F)) > 0x0F)
	c.setFlag(flagC, ((oldSP&0xFF)+(uint16(n)&0xFF)) > 0xFF)

	c.sp = uint16(result)

	return 16

	// log.Info("add_sp_dd\tADD SP, %02X", n)
}

// ld_hl_sp_dd
//
// Clock cycles: 12
func (c *cpu) ld_hl_sp_dd() uint8 {
	n := int8(gb.bus.read(c.pc))
	c.pc += 1

	oldSP := c.sp
	result := int32(oldSP) + int32(n)

	c.setFlag(flagZ, false)
	c.setFlag(flagN, false)
	c.setFlag(flagH, ((oldSP&0x0F)+(uint16(n)&0x0F)) > 0x0F)
	c.setFlag(flagC, ((oldSP&0xFF)+(uint16(n)&0xFF)) > 0xFF)

	c.writeRegister16(reg_hl, uint16(result))
	return 12
	// log.Info("ld_hl_sp_dd\tLD HL, SP+%02X", n)
}

/** Rotate and shift instructions **/

// rlca: Rotate register A left.
//
// Clock cycles: 4
func (c *cpu) rlca() uint8 {
	a := c.readRegister(reg_a)
	carry := a >> 7

	result := (a << 1) | carry

	c.setFlag(flagZ, false)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, carry != 0)

	c.writeRegister(reg_a, result)
	return 4
	// log.Info("rlca\tRLCA A[%02X]", result)
}

// rla
//
// Clock cycles: 4
func (c *cpu) rla() uint8 {
	a := c.readRegister(reg_a)
	carry := c.getFlag(flagC)

	result := (a << 1) | boolToUint8(carry)

	c.setFlag(flagZ, false)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, a>>7 != 0)

	c.writeRegister(reg_a, result)
	return 4
	// log.Info("rla\tRLA A[%02X]", result)
}

// rrca
//
// Clock cycles: 4
func (c *cpu) rrca() uint8 {
	a := c.readRegister(reg_a)
	carry := a & 0x01

	result := (a >> 1) | (carry << 7)

	c.setFlag(flagZ, false)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, carry != 0)

	c.writeRegister(reg_a, result)
	return 4
	// log.Info("rrca\tRRCA A[%02X]", result)
}

// rra
//
// Clock cycles: 4
func (c *cpu) rra() uint8 {
	a := c.readRegister(reg_a)
	carry := c.getFlag(flagC)

	result := (a >> 1) | (boolToUint8(carry) << 7)

	c.setFlag(flagZ, false)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, (a&0x01) != 0)

	c.writeRegister(reg_a, result)
	return 4
	// log.Info("rra\tRRA A[%02X]", result)
}

// rlc_r
//
// Clock cycles: 8
func (c *cpu) rlc_r(r register) uint8 {
	value := c.readRegister(r)
	carry := value >> 7

	result := (value << 1) | carry

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, carry != 0)

	c.writeRegister(r, result)
	return 8
	// log.Info("rlc_%s\tRLC %s[%02X]", r.name(), r.name(), result)
}

// rlc_HL
//
// Clock cycles: 16
func (c *cpu) rlc_HL() uint8 {
	value := gb.bus.read(c.readRegister16(reg_hl))
	carry := value >> 7

	result := (value << 1) | carry

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, carry != 0)

	gb.bus.write(c.readRegister16(reg_hl), result)
	return 16
	// log.Info("rlc_HL\tRLC (HL)[%02X]", result)
}

// rl_r
//
// Clock cycles: 8
func (c *cpu) rl_r(r register) uint8 {
	value := c.readRegister(r)
	carry := c.getFlag(flagC)

	result := (value << 1) | boolToUint8(carry)

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, value>>7 != 0)

	c.writeRegister(r, result)
	return 8
	// log.Info("rl_%s\tRL %s[%02X]", r.name(), r.name(), result)
}

// rl_HL
//
// Clock cycles: 16
func (c *cpu) rl_HL() uint8 {
	value := gb.bus.read(c.readRegister16(reg_hl))
	carry := c.getFlag(flagC)

	result := (value << 1) | boolToUint8(carry)

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, value>>7 != 0)

	gb.bus.write(c.readRegister16(reg_hl), result)
	return 16
	// log.Info("rl_HL\tRL (HL)[%02X]", result)
}

// rrc_r
//
// Clock cycles: 8
func (c *cpu) rrc_r(r register) uint8 {
	value := c.readRegister(r)
	carry := value & 0x01

	result := (value >> 1) | (carry << 7)

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, carry != 0)

	c.writeRegister(r, result)
	return 8
	// log.Info("rrc_%s\tRRC %s[%02X]", r.name(), r.name(), result)
}

// rrc_HL
//
// Clock cycles: 16
func (c *cpu) rrc_HL() uint8 {
	value := gb.bus.read(c.readRegister16(reg_hl))
	carry := value & 0x01

	result := (value >> 1) | (carry << 7)

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, carry != 0)

	gb.bus.write(c.readRegister16(reg_hl), result)
	return 16
	// log.Info("rrc_HL\tRRC (HL)[%02X]", result)
}

// rr_r
//
// Clock cycles: 8
func (c *cpu) rr_r(r register) uint8 {
	value := c.readRegister(r)
	carry := c.getFlag(flagC)

	result := (value >> 1) | (boolToUint8(carry) << 7)

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, (value&0x01) != 0)

	c.writeRegister(r, result)
	return 8
	// log.Info("rr_%s\tRR %s[%02X]", r.name(), r.name(), result)
}

// rr_HL
//
// Clock cycles: 16
func (c *cpu) rr_HL() uint8 {
	value := gb.bus.read(c.readRegister16(reg_hl))
	carry := c.getFlag(flagC)

	result := (value >> 1) | (boolToUint8(carry) << 7)

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, (value&0x01) != 0)

	gb.bus.write(c.readRegister16(reg_hl), result)
	return 16
	// log.Info("rr_HL\tRR (HL)[%02X]", result)
}

// sla_r
//
// Clock cycles: 8
func (c *cpu) sla_r(r register) uint8 {
	value := c.readRegister(r)

	result := value << 1

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, value>>7 != 0)

	c.writeRegister(r, result)
	return 8
	// log.Info("sla_%s\tSLA %s[%02X]", r.name(), r.name(), result)
}

// sla_HL
//
// Clock cycles: 16
func (c *cpu) sla_HL() uint8 {
	value := gb.bus.read(c.readRegister16(reg_hl))

	result := value << 1

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, value>>7 != 0)

	gb.bus.write(c.readRegister16(reg_hl), result)
	return 16
	// log.Info("sla_HL\tSLA (HL)[%02X]", result)
}

// swap_r
//
// Clock cycles: 8
func (c *cpu) swap_r(r register) uint8 {
	value := c.readRegister(r)

	result := ((value & 0x0F) << 4) | ((value & 0xF0) >> 4)

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, false)

	c.writeRegister(r, result)
	return 8
	// log.Info("swap_%s\tSWAP %s[%02X]", r.name(), r.name(), result)
}

// swap_HL
//
// Clock cycles: 16
func (c *cpu) swap_HL() uint8 {
	value := gb.bus.read(c.readRegister16(reg_hl))

	result := ((value & 0x0F) << 4) | ((value & 0xF0) >> 4)

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, false)

	gb.bus.write(c.readRegister16(reg_hl), result)
	return 16
	// log.Info("swap_HL\tSWAP (HL)[%02X]", result)
}

// sra_r
//
// Clock cycles: 8
func (c *cpu) sra_r(r register) uint8 {
	value := c.readRegister(r)
	msb := value & 0x80

	result := (value >> 1) | msb

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, value&0x01 != 0)

	c.writeRegister(r, result)
	return 8
	// log.Info("sra_%s\tSRA %s[%02X]", r.name(), r.name(), result)
}

// sra_HL
//
// Clock cycles: 16
func (c *cpu) sra_HL() uint8 {
	value := gb.bus.read(c.readRegister16(reg_hl))
	msb := value & 0x80

	result := (value >> 1) | msb

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, value&0x01 != 0)

	gb.bus.write(c.readRegister16(reg_hl), result)
	return 16
	// log.Info("sra_HL\tSRA (HL)[%02X]", result)
}

// srl_r
//
// Clock cycles: 8
func (c *cpu) srl_r(r register) uint8 {
	value := c.readRegister(r)

	result := value >> 1

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, value&0x01 != 0)

	c.writeRegister(r, result)
	return 8
	// log.Info("srl_%s\tSRL %s[%02X]", r.name(), r.name(), result)
}

// srl_HL
//
// Clock cycles: 16
func (c *cpu) srl_HL() uint8 {
	value := gb.bus.read(c.readRegister16(reg_hl))

	result := value >> 1

	c.setFlag(flagZ, result == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, value&0x01 != 0)

	gb.bus.write(c.readRegister16(reg_hl), result)
	return 16
	// log.Info("srl_HL\tSRL (HL)[%02X]", result)
}

/** Single-bit Operation instructions **/

// bit_n_r
//
// Clock cycles: 8
func (c *cpu) bit_n_r(bit int, r register) uint8 {
	value := c.readRegister(r)

	c.setFlag(flagZ, (value>>bit)&1 == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, true)

	return 8
	// log.Info("bit_%d_%s\tBIT %d, %s", bit, r.name(), bit, r.name())
}

// bit_n_HL
//
// Clock cycles: 12
func (c *cpu) bit_n_HL(bit int) uint8 {
	value := gb.bus.read(c.readRegister16(reg_hl))

	c.setFlag(flagZ, (value>>bit)&1 == 0)
	c.setFlag(flagN, false)
	c.setFlag(flagH, true)

	return 12
	// log.Info("bit_%d_HL\tBIT %d, (HL)", bit, bit)
}

// set_n_r
//
// Clock cycles: 8
func (c *cpu) set_n_r(bit int, r register) uint8 {
	value := c.readRegister(r)
	result := value | (1 << bit)

	c.writeRegister(r, result)
	return 8
	// log.Info("set_%d_%s\tSET %d, %s", bit, r.name(), bit, r.name())
}

// set_n_HL
//
// Clock cycles: 16
func (c *cpu) set_n_HL(bit int) uint8 {
	value := gb.bus.read(c.readRegister16(reg_hl))
	result := value | (1 << bit)

	gb.bus.write(c.readRegister16(reg_hl), result)
	return 16
	// log.Info("set_%d_HL\tSET %d, (HL)", bit, bit)
}

// res_n_r
//
// Clock cycles: 8
func (c *cpu) res_n_r(bit int, r register) uint8 {
	value := c.readRegister(r)
	result := value & ^(1 << bit)

	c.writeRegister(r, result)
	return 8
	// log.Info("res_%d_%s\tRES %d, %s", bit, r.name(), bit, r.name())
}

// res_n_HL
//
// Clock cycles: 16
func (c *cpu) res_n_HL(bit int) uint8 {
	value := gb.bus.read(c.readRegister16(reg_hl))
	result := value & ^(1 << bit)

	gb.bus.write(c.readRegister16(reg_hl), result)
	return 16
	// log.Info("res_%d_HL\tRES %d, (HL)", bit, bit)
}

/** CPU Control instructions **/

// ccf: Complement Carry Flag
//
// Clock cycles: 4
func (c *cpu) ccf() uint8 {
	carry := c.getFlag(flagC)
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, !carry)

	return 4
}

// scf: Set Carry Flag
//
// Clock cycles: 4
func (c *cpu) scf() uint8 {
	c.setFlag(flagN, false)
	c.setFlag(flagH, false)
	c.setFlag(flagC, true)

	return 4
}

// halt
//
// Clock cycles: N*4
func (c *cpu) halt() uint8 {
	// Implement halt logic based on your system's requirements
	// This instruction stops the CPU until an interrupt occurs
	// You might need to use sleep or other system-specific calls
	// to properly implement the halt instruction
	// N*4 cycles where N is the number of cycles the CPU is halted
	c.halted = true
	return 4
}

// stop
//
// Clock cycles: ?
func (c *cpu) stop() uint8 {
	// Implement stop logic based on your system's requirements
	// This instruction puts the CPU into a low-power standby mode
	// You might need to use sleep or other system-specific calls
	// to properly implement the stop instruction
	panic("stop")
}

// di: Disable Interrupts
//
// Clock cycles: 4
func (c *cpu) di() uint8 {
	// Implement interrupt disable logic (IME = 0)
	// Disable interrupts so that they cannot occur
	c.ime = false

	return 4
}

// ei: Enable Interrupts
//
// Clock cycles: 4
func (c *cpu) ei() uint8 {
	// Implement interrupt enable logic (IME = 1)
	// Enable interrupts so that they can occur
	c.ime = true

	return 4
}

/** Jump instructions **/

// jp_nn: Jump to nn
//
// Clock cycles: 16
func (c *cpu) jp_nn() uint8 {
	address := gb.bus.read16(c.pc)
	c.pc = address
	return 16
	// log.Info("jp_nn\tJP %04X", address)
}

// jp_HL
//
// Clock cycles: 4
func (c *cpu) jp_HL() uint8 {
	address := c.readRegister16(reg_hl)
	c.pc = address
	return 4
	// log.Info("jp_HL\tJP HL[%04X]", address)
}

// jp_f_nn
//
// Clock cycles: 12(false)/16(true)
func (c *cpu) jp_f_nn(condition condition) uint8 {
	address := gb.bus.read16(c.pc)
	c.pc += 2

	cycles := uint8(0)
	if c.evalCond(condition) {
		c.pc = address
		cycles = 16
	} else {
		cycles = 12
	}
	// log.Info("jp_%s_nn\tJP %s, %04X", condition.toString(), condition.toString(), address)
	return cycles
}

// jr_PC_dd: Relative Jump
//
// Clock cycles: 12
func (c *cpu) jr_PC_dd() uint8 {
	offset := gb.bus.read(c.pc)
	c.pc += 1
	newPC := uint16(int(c.pc) + int(int8(offset)))
	c.pc = newPC
	return 12
	// log.Info("jr_PC_dd\tJR %04X", newPC)
}

// jr_f_PC_dd: Conditional Relative Jump
//
// Clock cycles: 8(false)/12(true)
func (c *cpu) jr_f_PC_dd(condition condition) uint8 {
	offset := gb.bus.read(c.pc)
	c.pc += 1
	var cycles uint8
	if c.evalCond(condition) {
		newPC := uint16(int(c.pc) + int(int8(offset)))
		c.pc = newPC
		cycles = 12
	} else {
		cycles = 8
	}
	// log.Info("jr_%s_PC_dd\tJR %s, %04X", condition.toString(), condition.toString(), c.pc)
	return cycles
}

// call_nn: Call nn
//
// Clock cycles: 24
func (c *cpu) call_nn() uint8 {
	address := gb.bus.read16(c.pc)
	c.pc += 2
	c.push16(c.pc)
	c.pc = address
	return 24
	// log.Info("call_nn\tCALL %04X", address)
}

// call_f_nn: Conditional Call
//
// Clock cycles: 12(false)/24(true)
func (c *cpu) call_f_nn(condition condition) uint8 {
	address := gb.bus.read16(c.pc)
	c.pc += 2
	var cycles uint8
	if c.evalCond(condition) {
		c.push16(c.pc)
		c.pc = address
		cycles = 24
	} else {
		cycles = 12
	}
	// log.Info("call_%s_nn\tCALL %s, %04X", condition.toString(), condition.toString(), address)
	return cycles
}

// ret: Return
//
// Clock cycles: 16
func (c *cpu) ret() uint8 {
	address := c.pop16()
	c.pc = address
	return 16
}

// ret_f: Conditional Return
//
// Clock cycles: 8(false)/20(true)
func (c *cpu) ret_f(condition condition) uint8 {
	var cycles uint8
	if c.evalCond(condition) {
		address := c.pop16()
		c.pc = address
		cycles = 20
	} else {
		cycles = 8
	}
	// log.Info("ret_%s\tRET %s", condition.toString(), condition.toString())
	return cycles
}

// reti: Return and Enable Interrupts
//
// Clock cycles: 16
func (c *cpu) reti() uint8 {
	c.ime = true
	return c.ret() // ret already has 16 clock cycles
	// log.Info("reti\tRETI")
}

// rst_n: Restart
//
// Clock cycles: 16
func (c *cpu) rst_n(addr uint16) uint8 {
	c.push16(c.pc)
	c.pc = addr
	return 16
	// log.Info("rst_%02X\tRST %02X", addr, addr)
}
