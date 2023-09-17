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
			"../../roms/01-special.gb",
			2000,
			"01-specialPassed",
		},
		{
			"02-interrupts",
			"../../roms/02-interrupts.gb",
			2000,
			"02-interruptsPassed",
		},
		{
			"03-op sp,hl.gb",
			"../../roms/03-op sp,hl.gb",
			2000,
			"03-op sp,hlPassed",
		},
		{
			"04-op r,imm.gb",
			"../../roms/04-op r,imm.gb",
			2000,
			"04-op r,immPassed",
		},
		{
			"05-op rp.gb",
			"../../roms/05-op rp.gb",
			2000,
			"05-op rpPassed",
		},
		{
			"06-ld r,r.gb",
			"../../roms/06-ld r,r.gb",
			2000,
			"06-ld r,rPassed",
		},
		{
			"07-jr,jp,call,ret,rst.gb",
			"../../roms/07-jr,jp,call,ret,rst.gb",
			2000,
			"07-jr,jp,call,ret,rstPassed",
		},
		{
			"08-misc instrs.gb",
			"../../roms/08-misc instrs.gb",
			2000,
			"08-misc instrsPassed",
		},
		{
			"09-op r,r.gb",
			"../../roms/09-op r,r.gb",
			2000,
			"09-op r,rPassed",
		},
		{
			"10-bit ops.gb",
			"../../roms/10-bit ops.gb",
			2000,
			"10-bit opsPassed",
		},
		{
			"11-op a,(hl).gb",
			"../../roms/11-op a,(hl).gb",
			2000,
			"11-op a,(hl)Passed",
		},
		{
			"instr_timing.gb",
			"../../roms/instr_timing.gb",
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
