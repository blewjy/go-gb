package main

import (
	"fmt"
	"time"

	"github.com/blewjy/fire-gb/gb"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type game struct {
	ticks uint64
	start time.Time
}

func (g *game) Update() error {
	g.ticks++
	g.updateTitle()

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		gb.EnableDebug()
	}

	gb.Update()

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *game) updateTitle() {
	realtps := float64(g.ticks) / float64(time.Since(g.start).Seconds())
	ebiten.SetWindowTitle(fmt.Sprintf("fire-gb | Updates/s: %.2f Ticks: %d TPS: %.2f FPS: %.2f", realtps, g.ticks, ebiten.ActualTPS(), ebiten.ActualFPS()))
}
