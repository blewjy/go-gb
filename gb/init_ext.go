package gb

type TestGameboy struct {
	GB *gameboy
}

type TestGameboyStateCPU struct {
	A  uint8
	B  uint8
	C  uint8
	D  uint8
	E  uint8
	F  uint8
	H  uint8
	L  uint8
	PC uint16
	SP uint16
}

type TestGameboyState struct {
	CPU TestGameboyStateCPU
	RAM map[uint16]uint8
	ROM []uint8
}

func InitTestGameboy(state TestGameboyState) *TestGameboy {
	tgb := &gameboy{
		rom: state.ROM,

		bus:   newBus(),
		cpu:   newCpu(),
		timer: newTimer(),
		ram:   newRam(),
		cart:  newCart(state.ROM),
	}

	tgb.cpu.a = state.CPU.A
	tgb.cpu.b = state.CPU.B
	tgb.cpu.c = state.CPU.C
	tgb.cpu.d = state.CPU.D
	tgb.cpu.e = state.CPU.E
	tgb.cpu.f = state.CPU.F
	tgb.cpu.h = state.CPU.H
	tgb.cpu.l = state.CPU.L
	tgb.cpu.pc = state.CPU.PC
	tgb.cpu.sp = state.CPU.SP

	for addr, val := range state.RAM {
		tgb.bus.write(addr, val)
	}

	return &TestGameboy{GB: tgb}
}

func (tgb *TestGameboy) ExportState() TestGameboyState {
	convertRamToMap := func(ram []uint8) map[uint16]uint8 {
		m := map[uint16]uint8{}
		for addr, val := range ram {
			m[uint16(addr)] = val
		}
		return m
	}
	return TestGameboyState{
		CPU: TestGameboyStateCPU{
			A:  tgb.GB.cpu.a,
			B:  tgb.GB.cpu.b,
			C:  tgb.GB.cpu.c,
			D:  tgb.GB.cpu.d,
			E:  tgb.GB.cpu.e,
			F:  tgb.GB.cpu.f,
			H:  tgb.GB.cpu.h,
			L:  tgb.GB.cpu.l,
			PC: tgb.GB.cpu.pc,
			SP: tgb.GB.cpu.sp,
		},
		RAM: convertRamToMap(tgb.GB.ram.memory),
	}
}

func (tgb *TestGameboy) StepCPU() {
	tgb.GB.cpu.step()
}
