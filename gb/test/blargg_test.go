package test

import (
	"os"
	"strings"
	"testing"

	"github.com/blewjy/fire-gb/gb"
)

func TestBlargg(t *testing.T) {
	tests := []struct {
		name    string
		romName string
		clocks  int
		want    string
	}{
		{
			"01-special",
			"../../roms/blargg/01-special.gb",
			2000,
			"01-specialPassed",
		},
		{
			"02-interrupts",
			"../../roms/blargg/02-interrupts.gb",
			2000,
			"02-interruptsPassed",
		},
		{
			"03-op sp,hl.gb",
			"../../roms/blargg/03-op sp,hl.gb",
			2000,
			"03-op sp,hlPassed",
		},
		{
			"04-op r,imm.gb",
			"../../roms/blargg/04-op r,imm.gb",
			2000,
			"04-op r,immPassed",
		},
		{
			"05-op rp.gb",
			"../../roms/blargg/05-op rp.gb",
			2000,
			"05-op rpPassed",
		},
		{
			"06-ld r,r.gb",
			"../../roms/blargg/06-ld r,r.gb",
			2000,
			"06-ld r,rPassed",
		},
		{
			"07-jr,jp,call,ret,rst.gb",
			"../../roms/blargg/07-jr,jp,call,ret,rst.gb",
			2000,
			"07-jr,jp,call,ret,rstPassed",
		},
		{
			"08-misc instrs.gb",
			"../../roms/blargg/08-misc instrs.gb",
			2000,
			"08-misc instrsPassed",
		},
		{
			"09-op r,r.gb",
			"../../roms/blargg/09-op r,r.gb",
			2000,
			"09-op r,rPassed",
		},
		{
			"10-bit ops.gb",
			"../../roms/blargg/10-bit ops.gb",
			2000,
			"10-bit opsPassed",
		},
		{
			"11-op a,(hl).gb",
			"../../roms/blargg/11-op a,(hl).gb",
			2000,
			"11-op a,(hl)Passed",
		},
		{
			"instr_timing.gb",
			"../../roms/blargg/instr_timing.gb",
			2000,
			"instr_timingPassed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			romBytes, err := os.ReadFile(tt.romName)
			if err != nil {
				panic(err)
			}
			gameboy := gb.InitWithoutDisplay(romBytes)
			gameboy.SetHeadless(true)
			gameboy.SetTestMode(true)
			for i := 0; i < tt.clocks; i++ {
				gameboy.Update()
				got := strings.ReplaceAll(gameboy.GetDebug(), "\n", "")
				if got == tt.want {
					break
				}
			}
			got := strings.ReplaceAll(gameboy.GetDebug(), "\n", "")
			if got != tt.want {
				t.Errorf("%v: got = %v, want = %v", tt.name, got, tt.want)
			}
		})
	}
}
