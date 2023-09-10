package gb

type timer struct {
	div  uint16
	tima uint8
	tma  uint8
	tac  uint8
}

func newTimer() *timer {
	return &timer{
		div: 0xABCC,
	}
}

func (t *timer) clock() {
	prevDiv := t.div
	t.div += 1

	timaUpdate := false

	switch t.tac & 0b11 {
	case 0b00:
		timaUpdate = (prevDiv&(1<<9) > 0) && (t.div&(1<<9) == 0)
		break
	case 0b01:
		timaUpdate = (prevDiv&(1<<3) > 0) && (t.div&(1<<3) == 0)
		break
	case 0b10:
		timaUpdate = (prevDiv&(1<<5) > 0) && (t.div&(1<<5) == 0)
		break
	case 0b11:
		timaUpdate = (prevDiv&(1<<7) > 0) && (t.div&(1<<7) == 0)
		break
	}

	if timaUpdate && t.tac&4 > 0 {
		t.tima += 1
		if t.tima == 0xFF {
			t.tima = t.tma
			gb.cpu.reqInterrupt(cpuInterruptTimer)
		}
	}
}

func (t *timer) read(addr uint16) uint8 {
	switch addr {
	case 0xFF04:
		return uint8(t.div >> 8)
	case 0xFF05:
		return t.tima
	case 0xFF06:
		return t.tma
	case 0xFF07:
		return t.tac
	default:
		panic("out of range")
	}
}

func (t *timer) write(addr uint16, value uint8) {
	switch addr {
	case 0xFF04:
		t.div = 0
	case 0xFF05:
		t.tima = value
	case 0xFF06:
		t.tma = value
	case 0xFF07:
		t.tac = value
	default:
		panic("out of range")
	}
}
