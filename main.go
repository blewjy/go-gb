package main

import (
	"github.com/blewjy/fire-gb/gb"
	"github.com/hajimehoshi/ebiten/v2"
	"os"
	"time"
)

func main() {
	ebiten.SetWindowTitle("go-gameboy")
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// romBytes, err := os.ReadFile("roms/dmg-acid2.gb")
	// romBytes, err := os.ReadFile("roms/01-special.gb")
	romBytes, err := os.ReadFile("roms/02-interrupts.gb")
	// romBytes, err := os.ReadFile("roms/03-op sp,hl.gb")
	// romBytes, err := os.ReadFile("roms/04-op r,imm.gb")
	// romBytes, err := os.ReadFile("roms/05-op rp.gb")
	// romBytes, err := os.ReadFile("roms/06-ld r,r.gb")
	// romBytes, err := os.ReadFile("roms/07-jr,jp,call,ret,rst.gb")
	// romBytes, err := os.ReadFile("roms/08-misc instrs.gb")
	// romBytes, err := os.ReadFile("roms/09-op r,r.gb")
	// romBytes, err := os.ReadFile("roms/10-bit ops.gb")
	// romBytes, err := os.ReadFile("roms/11-op a,(hl).gb")
	if err != nil {
		panic(err)
	}

	gb.InitWithROM(romBytes)

	g := &game{
		start: time.Now(),
	}
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
