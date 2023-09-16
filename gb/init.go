package gb

type Gameboy struct {
	rom []byte

	bus   *bus
	cpu   *cpu
	timer *timer
	ram   *ram
	cart  *cart

	shouldDebug bool
	debug       bool
}

func InitWithROM(rom []byte) *Gameboy {
	gb := &Gameboy{
		rom: rom,
	}

	gb.bus = newBus(gb)
	gb.cpu = newCpu(gb)
	gb.timer = newTimer(gb)
	gb.ram = newRam()
	gb.cart = newCart(rom)

	return gb
}

// Update must be called at 60Hz
func (gb *Gameboy) Update() {
	currCycles := uint64(0)
	for currCycles < cpuFreq/60 {
		if gb.shouldDebug {
			gb.debug = true
			gb.shouldDebug = false
		}

		cycles := gb.cpu.step()

		for i := uint8(0); i < cycles; i++ {
			gb.timer.clock()
		}

		currCycles += uint64(cycles)

		gb.debug = false
	}
}

func (gb *Gameboy) StepCPU() {
	gb.cpu.step()
}

func (gb *Gameboy) EnableDebug() {
	gb.shouldDebug = true
}

func (gb *Gameboy) GetDebug() string {
	return gb.cpu.debugMsg
}
