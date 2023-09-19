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

	//romBytes, err := os.ReadFile("roms/dmg-acid2.gb")
	romBytes, err := os.ReadFile("roms/drmario.gb")
	//romBytes, err := os.ReadFile("roms/sml.gb")
	//romBytes, err := os.ReadFile("roms/tetris.gb")
	if err != nil {
		panic(err)
	}

	g := &game{
		gb: gb.Init(romBytes),

		start: time.Now(),

		screen: ebiten.NewImage(160, 144),
	}
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
