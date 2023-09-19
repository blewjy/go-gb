package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/blewjy/fire-gb/gb"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type game struct {
	gb *gb.Gameboy

	ticks uint64
	start time.Time

	screen *ebiten.Image
}

func (g *game) Update() error {
	g.ticks++
	g.updateTitle()
	g.handleInputs()

	g.gb.Update()

	// update the screen
	colorToBytes := func(c color.RGBA) []byte {
		return []byte{c.R, c.G, c.B, c.A}
	}

	var displayBytes []byte
	emuDisplay := g.gb.GetDisplay()
	for _, row := range emuDisplay {
		for _, pixel := range row {
			displayBytes = append(displayBytes, colorToBytes(pixel)...)
		}
	}
	g.screen.WritePixels(displayBytes)
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	op.GeoM.Scale(4, 4)
	screen.DrawImage(g.screen, op)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *game) updateTitle() {
	realtps := float64(g.ticks) / float64(time.Since(g.start).Seconds())
	ebiten.SetWindowTitle(fmt.Sprintf("fire-gb | Updates/s: %.2f Ticks: %d TPS: %.2f FPS: %.2f", realtps, g.ticks, ebiten.ActualTPS(), ebiten.ActualFPS()))
}

var joypadInputMap = map[ebiten.Key]gb.JoypadButton{
	ebiten.KeyEnter:     gb.JoypadButtonStart,
	ebiten.KeyBackspace: gb.JoypadButtonSelect,
	ebiten.KeyZ:         gb.JoypadButtonA,
	ebiten.KeyX:         gb.JoypadButtonB,
	ebiten.KeyDown:      gb.JoypadButtonDown,
	ebiten.KeyUp:        gb.JoypadButtonUp,
	ebiten.KeyLeft:      gb.JoypadButtonLeft,
	ebiten.KeyRight:     gb.JoypadButtonRight,
}

func (g *game) handleInputs() {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.gb.EnableDebug()
	}

	for key, button := range joypadInputMap {
		if inpututil.IsKeyJustPressed(key) {
			g.gb.SetButtonPressed(button)
		}
		if inpututil.IsKeyJustReleased(key) {
			g.gb.SetButtonReleased(button)
		}
	}
}
