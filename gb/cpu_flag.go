package gb

type flag uint8

const (
	flagZ flag = 7 // Zero flag
	flagN flag = 6 // Subtract flag
	flagH flag = 5 // Half carry flag
	flagC flag = 4 // Carry flag
)

func (c *cpu) getFlag(f flag) bool {
	return (c.f & (1 << uint8(f))) != 0
}

func (c *cpu) setFlag(f flag, value bool) {
	if value {
		c.f |= (1 << uint8(f))
	} else {
		c.f &^= (1 << uint8(f))
	}
}

type condition int

const (
	conditionNZ condition = iota
	conditionZ
	conditionNC
	conditionC
)

func (c *cpu) evalCond(condition condition) bool {
	var result bool
	switch condition {
	case conditionNZ:
		result = c.getFlag(flagZ) == false
	case conditionZ:
		result = c.getFlag(flagZ) == true
	case conditionNC:
		result = c.getFlag(flagC) == false
	case conditionC:
		result = c.getFlag(flagC) == true
	default:
		panic("invalid condition")
	}
	return result
}

func (c condition) toString() string {
	switch c {
	case conditionNZ:
		return "NZ"
	case conditionZ:
		return "Z"
	case conditionNC:
		return "NC"
	case conditionC:
		return "C"
	default:
		panic("invalid condition")
	}
}
