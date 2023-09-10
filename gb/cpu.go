package gb

import "fmt"

const (
	cpuFreq = 4194304
)

type cpu struct {
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

func newCpu() *cpu {
	return &cpu{
		pc: 0x0100,
		sp: 0xFFFE,

		cycles: 0,
	}
}

func (c *cpu) step() uint8 {
	opcode := gb.bus.read(c.pc)
	// fmt.Printf("current PC: %04X, opcode: %02X\n", c.pc, opcode)
	c.pc += 1
	cycles := c.execute(opcode)

	c.cycles += uint64(cycles)

	c.serialDebug()
	if gb.debug && c.debugMsg != "" {
		fmt.Println(c.debugMsg)
	}

	return cycles
}

func (c *cpu) execute(opcode uint8) uint8 {
	return opcodeToInst[opcode](c)
}

func (c *cpu) serialDebug() {
	if gb.bus.read(0xFF02) == 0x81 {
		v := gb.bus.read(0xFF01)
		c.debugMsg += string(v)
		gb.bus.write(0xFF02, 0x01)
	}
}
