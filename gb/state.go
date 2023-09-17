package gb

type CPUState struct {
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

type State struct {
	CPU CPUState
	RAM map[uint16]uint8
	ROM []uint8
}

func (gb *Gameboy) SetState(state State) {
	gb.cpu.a = state.CPU.A
	gb.cpu.b = state.CPU.B
	gb.cpu.c = state.CPU.C
	gb.cpu.d = state.CPU.D
	gb.cpu.e = state.CPU.E
	gb.cpu.f = state.CPU.F
	gb.cpu.h = state.CPU.H
	gb.cpu.l = state.CPU.L
	gb.cpu.pc = state.CPU.PC
	gb.cpu.sp = state.CPU.SP

	for addr, val := range state.RAM {
		gb.bus.write(addr, val)
	}
}

func (gb *Gameboy) ExportStateWithAddresses(addresses []uint16) State {
	getBusData := func(addresses []uint16) map[uint16]uint8 {
		m := map[uint16]uint8{}
		for _, addr := range addresses {
			m[addr] = gb.bus.read(addr)
		}
		return m
	}
	return State{
		CPU: CPUState{
			A:  gb.cpu.a,
			B:  gb.cpu.b,
			C:  gb.cpu.c,
			D:  gb.cpu.d,
			E:  gb.cpu.e,
			F:  gb.cpu.f,
			H:  gb.cpu.h,
			L:  gb.cpu.l,
			PC: gb.cpu.pc,
			SP: gb.cpu.sp,
		},
		RAM: getBusData(addresses),
	}
}
