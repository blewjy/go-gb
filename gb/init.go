package gb

// Global gb master struct
var gb *gameboy

type gameboy struct {
	rom []byte

	bus   *bus
	cpu   *cpu
	timer *timer
	ram   *ram
	cart  *cart

	shouldDebug bool
	debug       bool
}

func InitWithROM(rom []byte) {
	gb = &gameboy{
		rom: rom,

		bus:   newBus(),
		cpu:   newCpu(),
		timer: newTimer(),
		ram:   newRam(),
		cart:  newCart(rom),
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

		for i := uint8(0); i < cycles; i++ {
			gb.timer.clock()
		}

		currCycles += uint64(cycles)

		gb.debug = false
	}
}

func EnableDebug() {
	gb.shouldDebug = true
}

func GetDebug() string {
	return gb.cpu.debugMsg
}
