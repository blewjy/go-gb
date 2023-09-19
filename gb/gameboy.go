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
	lcd   *lcd
	ppu   *ppu
	dma   *dma

	display [][]color.RGBA

	shouldDebug bool
	debug       bool
	headless    bool
	testMode    bool
}

func Init(rom []byte) *Gameboy {
	gb := InitWithoutDisplay(rom)

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
	gb.lcd = newLcd(gb)
	gb.ppu = newPpu(gb)
	gb.dma = newDma(gb)

	//gb.bus.write(0xff10, 0x80)
	//gb.bus.write(0xff11, 0xbf)
	//gb.bus.write(0xff12, 0xf3)
	//gb.bus.write(0xff14, 0xbf)
	//gb.bus.write(0xff17, 0x3f)
	//gb.bus.write(0xff19, 0xbf)
	//gb.bus.write(0xff1a, 0x7f)
	//gb.bus.write(0xff1b, 0xff)
	//gb.bus.write(0xff1c, 0x9f)
	//gb.bus.write(0xff1e, 0xbf)
	//gb.bus.write(0xff20, 0xff)
	//gb.bus.write(0xff23, 0xbf)
	//gb.bus.write(0xff24, 0x77)
	//gb.bus.write(0xff25, 0xf3)
	//gb.bus.write(0xff26, 0xf1)
	//
	//gb.bus.write(0xff40, 0x91)
	//gb.bus.write(0xff47, 0xfc)
	//gb.bus.write(0xff48, 0xff)
	//gb.bus.write(0xff49, 0xff)

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
			gb.dma.tick()
			gb.timer.clock()
			gb.ppu.tick()
		}

		currCycles += uint64(cycles)

		gb.debug = false

	}

	if !gb.headless && len(gb.ppu.framePxReady) == 23040 {
		c := 0
		for i := 0; i < 144; i++ {
			for j := 0; j < 160; j++ {
				px := gb.ppu.framePxReady[c]
				gb.display[i][j] = colorMap[colorChoice][px.color]
				c++
			}
		}
		//break
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

func (gb *Gameboy) SetHeadless(b bool) {
	gb.headless = b
}

func (gb *Gameboy) SetTestMode(b bool) {
	gb.testMode = b
}
