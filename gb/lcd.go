package gb

import "image/color"

type lcd struct {
	gb *Gameboy
}

func newLcd(gb *Gameboy) *lcd {
	l := &lcd{
		gb: gb,
	}
	l.set_STAT_modeflag(STAT_mode_searchoam)
	return l
}

func getBit(n uint8, pos uint8) uint8 {
	return (n >> pos) & 1
}

/*
LCDC uint8 // 0xFF40: LCDC (LCD control)
STAT uint8 // 0xFF41: STAT (LCD status)
SCY  uint8 // 0xFF42: SCY (Viewport Y-position)
SCX  uint8 // 0xFF43: SCX (Viewport X-position)
LY   uint8 // 0xFF44: LY (LCD Y-coordinate)
LYC  uint8 // 0xFF45: LYC (LY compare)
BGP  uint8 // 0xFF47: BGP (BG palette data)
OBP0 uint8 // 0xFF48: OBP0 (OBJ palette 0 data)
OBP1 uint8 // 0xFF49: OBP1 (OBJ palette 1 data)
WY   uint8 // 0xFF4A: WY (Window Y-position)
WX   uint8 // 0xFF4B: WX (Window X-position + 7)
*/

func (l *lcd) LCDC() uint8 {
	return l.gb.bus.read(0xFF40)
}

func (l *lcd) STAT() uint8 {
	return l.gb.bus.read(0xFF41)
}

func (l *lcd) SCY() uint8 {
	return l.gb.bus.read(0xFF42)
}

func (l *lcd) SCX() uint8 {
	return l.gb.bus.read(0xFF43)
}

func (l *lcd) LY() uint8 {
	return l.gb.bus.read(0xFF44)
}

func (l *lcd) incLY() {
	l.gb.bus.write(0xFF44, l.LY()+1)
}

func (l *lcd) resetLY() {
	l.gb.bus.write(0xFF44, 0)
}

func (l *lcd) LYC() uint8 {
	return l.gb.bus.read(0xFF45)
}

func (l *lcd) BGP() uint8 {
	return l.gb.bus.read(0xFF47)
}

func (l *lcd) OBP0() uint8 {
	return l.gb.bus.read(0xFF48)
}

func (l *lcd) OBP1() uint8 {
	return l.gb.bus.read(0xFF49)
}

func (l *lcd) WY() uint8 {
	return l.gb.bus.read(0xFF4A)
}

func (l *lcd) WX() uint8 {
	return l.gb.bus.read(0xFF4B)
}

/*
0xFF40 - LCDC
7	LCD and PPU enable	            0=Off, 1=On
6	Window tile map area	        0=9800-9BFF, 1=9C00-9FFF
5	Window enable	                0=Off, 1=On
4	BG and Window tile data area	0=8800-97FF, 1=8000-8FFF
3	BG tile map area	            0=9800-9BFF, 1=9C00-9FFF
2	OBJ size	                    0=8×8, 1=8×16
1	OBJ enable	                    0=Off, 1=On
0	BG and Window enable/priority	0=Off, 1=On
*/

func (l *lcd) LCDC_enable() bool {
	return getBit(l.LCDC(), 7) > 0
}

func (l *lcd) LCDC_windowtilemap() uint8 {
	return getBit(l.LCDC(), 6)
}

func (l *lcd) LCDC_windowenable() bool {
	return getBit(l.LCDC(), 5) > 0
}

func (l *lcd) LCDC_bgwindowtiledata() uint8 {
	return getBit(l.LCDC(), 4)
}

func (l *lcd) LCDC_bgtilemap() uint8 {
	return getBit(l.LCDC(), 3)
}

func (l *lcd) LCDC_objsize() uint8 {
	return getBit(l.LCDC(), 2)
}

func (l *lcd) LCDC_objenable() bool {
	return getBit(l.LCDC(), 1) > 0
}

func (l *lcd) LCDC_bgwindowenable() bool {
	return getBit(l.LCDC(), 0) > 0
}

/*
0xFF41 STAT
Bit 6 - LYC=LY STAT Interrupt source         (1=Enable) (Read/Write)
Bit 5 - Mode 2 OAM STAT Interrupt source     (1=Enable) (Read/Write)
Bit 4 - Mode 1 VBlank STAT Interrupt source  (1=Enable) (Read/Write)
Bit 3 - Mode 0 HBlank STAT Interrupt source  (1=Enable) (Read/Write)
Bit 2 - LYC=LY Flag                          (0=Different, 1=Equal) (Read Only)
Bit 1-0 - Mode Flag                          (Mode 0-3, see below) (Read Only)
          0: HBlank
          1: VBlank
          2: Searching OAM
          3: Transferring Data to LCD Controller
*/

func (l *lcd) STAT_lycinterruptsource() bool {
	return getBit(l.STAT(), 6) > 0
}

func (l *lcd) STAT_oaminterruptsource() bool {
	return getBit(l.STAT(), 5) > 0
}

func (l *lcd) STAT_vblankinterruptsource() bool {
	return getBit(l.STAT(), 4) > 0
}

func (l *lcd) STAT_hblankinterruptsource() bool {
	return getBit(l.STAT(), 3) > 0
}

func (l *lcd) get_STAT_lycflag() bool {
	return getBit(l.STAT(), 2) > 0
}

func (l *lcd) set_STAT_lycflag(equal bool) {
	stat := l.STAT()
	if equal {
		l.gb.bus.write(0xFF41, stat|0b00000100)
	} else {
		l.gb.bus.write(0xFF41, stat&0b11111011)
	}
}

type STAT_mode uint8

const (
	STAT_mode_hblank       STAT_mode = 0
	STAT_mode_vblank       STAT_mode = 1
	STAT_mode_searchoam    STAT_mode = 2
	STAT_mode_transferring STAT_mode = 3
)

func (l *lcd) get_STAT_modeflag() STAT_mode {
	return STAT_mode(l.STAT() & 0b11)
}

func (l *lcd) set_STAT_modeflag(mode STAT_mode) {
	l.gb.bus.write(0xFF41, (l.STAT()&0b11111100)|uint8(mode)&0b11)
}

/*
0xFF47 BGP, 0xFF48 OBP0, 0xFF49 OBP1
Bit 7-6 - Color for index 3
Bit 5-4 - Color for index 2
Bit 3-2 - Color for index 1
Bit 1-0 - Color for index 0
*/

func (l *lcd) BGP_getcolor(idx uint8) color.RGBA {
	colorIdx := l.BGP() >> (idx * 2) & 0b11
	return colorMap[colorChoice][colorIdx]
}

func (l *lcd) OBP0_getcolor(idx uint8) color.RGBA {
	colorIdx := l.OBP0() >> (idx * 2) & 0b11
	return colorMap[colorChoice][colorIdx]
}

func (l *lcd) OBP1_getcolor(idx uint8) color.RGBA {
	colorIdx := l.OBP1() >> (idx * 2) & 0b11
	return colorMap[colorChoice][colorIdx]
}
