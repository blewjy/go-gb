package gb

type joypad struct {
	start    bool
	sel      bool
	b        bool
	a        bool
	down     bool
	up       bool
	left     bool
	right    bool
	readMask uint8
}

func newJoypad() *joypad {
	return &joypad{}
}

func (j *joypad) read() uint8 {
	joypadConv := func(b bool) uint8 {
		if b {
			return 0
		}
		return 1
	}

	if j.readMask == 0b01 { // button keys
		return 0b11000000 | j.readMask<<4 | joypadConv(j.start)<<3 | joypadConv(j.sel)<<2 | joypadConv(j.b)<<1 | joypadConv(j.a)
	} else if j.readMask == 0b10 { // direction keys
		return 0b11000000 | j.readMask<<4 | joypadConv(j.down)<<3 | joypadConv(j.up)<<2 | joypadConv(j.left)<<1 | joypadConv(j.right)
	}
	return 0b11000000 | j.readMask<<4 | 0b1111
}

func (j *joypad) write(val uint8) {
	j.readMask = val >> 4
}
