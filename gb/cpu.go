package gb

import "fmt"

const (
	cpuFreq = 4194304
)

type cpu struct {
	gb *Gameboy

	cycles uint64

	a  uint8
	b  uint8
	c  uint8
	d  uint8
	e  uint8
	f  uint8
	h  uint8
	l  uint8
	sp uint16
	pc uint16

	halted bool
	ime    bool

	debugMsg string
}

func newCpu(gb *Gameboy) *cpu {
	return &cpu{
		gb: gb,

		pc: 0x0100,
		sp: 0xFFFE,

		cycles: 0,
	}
}

func (c *cpu) step() uint8 {
	cycles := uint8(4)
	if !c.halted {
		opcode := c.gb.bus.read(c.pc)
		c.pc += 1
		cycles = c.execute(opcode)
	} else {
		intFlags := c.gb.bus.read(0xFF0F)
		if intFlags > 0 {
			c.halted = false
		}
	}

	if c.ime {
		c.handleInterrupts()
	}

	c.serialDebug()
	if c.gb.debug && c.debugMsg != "" {
		fmt.Println(c.debugMsg)
	}

	c.cycles += uint64(cycles)
	return cycles
}

func (c *cpu) execute(opcode uint8) uint8 {
	return opcodeToInst[opcode](c)
}

func (c *cpu) serialDebug() {
	if c.gb.bus.read(0xFF02) == 0x81 {
		v := c.gb.bus.read(0xFF01)
		c.debugMsg += string(v)
		c.gb.bus.write(0xFF02, 0x01)
	}
}
