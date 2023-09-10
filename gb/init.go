package gb

// Global gb master struct
var gb *gameboy

type gameboy struct {
	rom []byte

	bus  *bus
	cpu  *cpu
	ram  *ram
	cart *cart

	shouldDebug bool
	debug       bool
}

func InitWithROM(rom []byte) {
	gb = &gameboy{
		rom: rom,

		bus:  newBus(),
		cpu:  newCpu(),
		ram:  newRam(),
		cart: newCart(rom),
	}
}

// Update must be called at 60Hz
func Update() {
	currCycles := uint64(0)
	for currCycles < cpuFreq/60 {
		if gb.shouldDebug {
			gb.debug = true
			gb.shouldDebug = false
		}

		cycles := gb.cpu.step()
		currCycles += uint64(cycles)

		gb.debug = false
	}
}

func EnableDebug() {
	gb.shouldDebug = true
}
