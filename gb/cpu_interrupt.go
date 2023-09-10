package gb

type cpuInterrupt uint8

const (
	cpuInterruptVBlank  cpuInterrupt = 0x01
	cpuInterruptLcdStat cpuInterrupt = 0x02
	cpuInterruptTimer   cpuInterrupt = 0x04
	cpuInterruptSerial  cpuInterrupt = 0x08
	cpuInterruptJoypad  cpuInterrupt = 0x10
)

func (c *cpu) reqInterrupt(i cpuInterrupt) {
	intFlags := gb.bus.read(0xFF0F) | uint8(i)
	gb.bus.write(0xFF0F, intFlags)
}

func (c *cpu) isBitSet(d uint8, bit uint8) bool {
	return d&bit > 0
}

func (c *cpu) commonInterrupt(interruptFlag uint8, interrupt cpuInterrupt, addr uint16) {
	c.push16(c.pc)
	c.pc = addr
	gb.bus.write(0xFF0F, interruptFlag&(^uint8(interrupt)))
	c.halted = false
	c.ime = false
}

func (c *cpu) handleInterrupts() {
	interruptEnable := gb.bus.read(0xFFFF)
	interruptFlag := gb.bus.read(0xFF0F)
	if c.isBitSet(interruptEnable, uint8(cpuInterruptVBlank)) && c.isBitSet(interruptFlag, uint8(cpuInterruptVBlank)) {
		c.commonInterrupt(interruptFlag, cpuInterruptVBlank, 0x0040)
	} else if c.isBitSet(interruptEnable, uint8(cpuInterruptLcdStat)) && c.isBitSet(interruptFlag, uint8(cpuInterruptLcdStat)) {
		c.commonInterrupt(interruptFlag, cpuInterruptLcdStat, 0x0048)
	} else if c.isBitSet(interruptEnable, uint8(cpuInterruptTimer)) && c.isBitSet(interruptFlag, uint8(cpuInterruptTimer)) {
		c.commonInterrupt(interruptFlag, cpuInterruptTimer, 0x0050)
	} else if c.isBitSet(interruptEnable, uint8(cpuInterruptSerial)) && c.isBitSet(interruptFlag, uint8(cpuInterruptSerial)) {
		c.commonInterrupt(interruptFlag, cpuInterruptSerial, 0x0058)
	} else if c.isBitSet(interruptEnable, uint8(cpuInterruptJoypad)) && c.isBitSet(interruptFlag, uint8(cpuInterruptJoypad)) {
		c.commonInterrupt(interruptFlag, cpuInterruptJoypad, 0x0060)
	}
}
