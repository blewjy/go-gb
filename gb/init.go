package gb

import (
	"image/color"
)

type Gameboy struct {
	rom []byte

	bus   *bus
	cpu   *cpu
	timer *timer
	ram   *ram
	cart  *cart

	display [][]color.RGBA

	shouldDebug bool
	debug       bool
}

func Init(rom []byte) *Gameboy {
	gb := &Gameboy{
		rom: rom,
	}

	gb.bus = newBus(gb)
	gb.cpu = newCpu(gb)
	gb.timer = newTimer(gb)
	gb.ram = newRam()
	gb.cart = newCart(rom)

	for i := 0; i < 144; i++ {
		gb.display = append(gb.display, []color.RGBA{})
		for j := 0; j < 160; j++ {
			gb.display[i] = append(gb.display[i], color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff})
		}
	}

	return gb
}

func InitWithoutDisplay(rom []byte) *Gameboy {
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

func (gb *Gameboy) GetDisplay() [][]color.RGBA {
	return gb.display
}
