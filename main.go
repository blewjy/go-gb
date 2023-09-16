package main

import (
	"github.com/blewjy/fire-gb/gb"
	"github.com/hajimehoshi/ebiten/v2"
	"os"
	"time"
)

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	romBytes, err := os.ReadFile("roms/dmg-acid2.gb")
	if err != nil {
		panic(err)
	}

	g := &game{
		gb: gb.InitWithROM(romBytes),
		
		start: time.Now(),
	}
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
