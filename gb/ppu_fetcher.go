package gb

type fetcherState uint8

const (
	fetcherStateReadTile fetcherState = iota
	fetcherStateReadLo
	fetcherStateReadHi
	fetcherStateSleep
)

type fetcher struct {
	ppu   *ppu
	state fetcherState

	prevInWindow     bool
	windowDidDisplay bool
	windowY          uint8

	fetcherX uint8

	targetTile uint8
	targetLo   uint8
	targetHi   uint8
}

func newFetcher(ppu *ppu) *fetcher {
	return &fetcher{
		ppu:   ppu,
		state: fetcherStateReadTile,
	}
}

func (f *fetcher) clock() {
	switch f.state {
	case fetcherStateReadTile:
		f.handleReadTile()
	case fetcherStateReadLo:
		f.handleReadLo()
	case fetcherStateReadHi:
		f.handleReadHi()
	case fetcherStateSleep:
		f.handleSleep()
	default:
		panic("invalid fetcher state")
	}
}

func (f *fetcher) inWindow() bool {
	WX, WY := f.ppu.gb.lcd.WX(), f.ppu.gb.lcd.WY()
	LX, LY := f.fetcherX, f.ppu.gb.lcd.LY()
	return f.ppu.gb.lcd.LCDC_windowenable() && LX+7 >= WX && LY >= WY
}

func (f *fetcher) handleReadTile() {
	LX, LY := f.fetcherX, f.ppu.gb.lcd.LY()
	SCX, SCY := f.ppu.gb.lcd.SCX(), f.ppu.gb.lcd.SCY()
	WX, windowLine := f.ppu.gb.lcd.WX(), f.windowY

	var bit10, x, y uint8
	if f.inWindow() {
		f.windowDidDisplay = true
		bit10 = f.ppu.gb.lcd.LCDC_windowtilemap()
		y = windowLine / 8
		x = (LX + 7 - WX) / 8
	} else {
		bit10 = f.ppu.gb.lcd.LCDC_bgtilemap()
		y = (LY + SCY) / 8
		x = (LX + SCX) / 8
	}

	addr := 0b10011<<11 | uint16(bit10)<<10 | uint16(y)<<5 | uint16(x)
	f.targetTile = f.ppu.gb.bus.read(addr)

	f.state = fetcherStateReadLo
}

func (f *fetcher) handleReadLo() {
	LY, SCY := f.ppu.gb.lcd.LY(), f.ppu.gb.lcd.SCY()

	var y uint8
	if f.inWindow() {
		y = f.windowY % 8
	} else {
		y = (LY + SCY) % 8
	}

	target_bit12 := 1
	lcdc_bit4 := f.ppu.gb.lcd.LCDC_bgwindowtiledata()
	if lcdc_bit4 == 1 || getBit(f.targetTile, 7) == 1 {
		target_bit12 = 0
	}

	addr := 0b100<<13 | uint16(target_bit12)<<12 | uint16(f.targetTile)<<4 | uint16(y)<<1
	f.targetLo = f.ppu.gb.bus.read(addr)

	f.state = fetcherStateReadHi
}

func (f *fetcher) handleReadHi() {
	LY, SCY := f.ppu.gb.lcd.LY(), f.ppu.gb.lcd.SCY()

	var y uint8
	if f.inWindow() {
		y = f.windowY % 8
	} else {
		y = (LY + SCY) % 8
	}

	target_bit12 := 1
	lcdc_bit4 := f.ppu.gb.lcd.LCDC_bgwindowtiledata()
	if lcdc_bit4 == 1 || getBit(f.targetTile, 7) == 1 {
		target_bit12 = 0
	}

	addr := 0b100<<13 | uint16(target_bit12)<<12 | uint16(f.targetTile)<<4 | uint16(y)<<1 | 0b1
	f.targetHi = f.ppu.gb.bus.read(addr)

	f.state = fetcherStateSleep
}

func (f *fetcher) handleSleep() {
	getPaletteBits := func(value, idx uint8) uint8 {
		mask := uint8(3) << uint(idx*2)
		result := (value & mask) >> uint(idx*2)
		return result
	}

	// not supposed to do anything, but we will just use this step to push the pixel into the fifo
	SCX := f.ppu.gb.lcd.SCX()
	discard := 0
	if f.fetcherX == 0 {
		discard = int(SCX % 8)
	}

	for i := 7 - discard; i >= 0; i-- {
		// In this step, we are using the fetcherX to determine whether we have hit the window tile or not.
		// Using fetcherX means that we do not need to clear to FIFO when we hit the window tile,
		// we can directly append the next 8 window tiles on top.
		if !f.prevInWindow && f.inWindow() {
			f.prevInWindow = true
			f.state = fetcherStateReadTile
			return
		}

		// In this step, we are using the fetcherX to determine whether we have hit any obj tile.
		// If we do hit some obj tile, we need to pause the bg/win fetching and switch over to fetch
		// the entire 8 pixels of the object first, and then enqueue it to the obj_fifo.
		// The position of the enqueue is very important. We need to fill up the obj_fifo with
		// buffer/transparent pixels first, up to the len(bg_fifo), then enqueue the 8 pixels of this obj.
		//
		// For each pixel we encounter here, we need to check if its transparent. If it is, then we
		// should continue to check the next object in the oamLineEntries, since the later objs should show
		// thru transparent pixels of earlier objects. If we scan through all the oamLineEntries and
		// we found matches, but we don't find any non-transparent pixels, then we still need to enqueue a
		// transparent pixel.
		objX := uint8(255) // definitely more than width of screen
		objColorIdx := uint8(0)
		objContest := false
		for _, obj := range f.ppu.oamLineEntries {
			if f.fetcherX >= obj.x-8 && f.fetcherX < obj.x {

				// if we already have a potential non-transparent pixel
				// then we first check if this obj has a smaller x (smaller x take priority)
				// so if this obj has a larger X or equal X, we can just skip it already
				if objColorIdx != 0 {
					if obj.x >= objX {
						continue
					}
				}

				y := (f.ppu.gb.lcd.LY() - obj.y) % 8

				// if y-flipped
				if getBit(obj.flags, 6) > 0 {
					y = 7 - y
				}

				addrLo := 0b1000<<12 | uint16(obj.tileIdx)<<4 | uint16(y)<<1
				tileLo := f.ppu.gb.bus.read(addrLo)
				addrHi := 0b1000<<12 | uint16(obj.tileIdx)<<4 | uint16(y)<<1 | 1
				tileHi := f.ppu.gb.bus.read(addrHi)

				bitLo := getBit(tileLo, uint8(obj.fetchedX))
				bitHi := getBit(tileHi, uint8(obj.fetchedX))

				// if X-flipped, we increment, else decrement
				if getBit(obj.flags, 5) > 0 {
					obj.fetchedX += 1
				} else {
					obj.fetchedX -= 1
				}

				colorIdx := bitHi<<1 | bitLo

				if colorIdx == 0 {
					continue
				}

				// check the palette (ignore 0)
				if colorIdx != 0 {
					if getBit(obj.flags, 4) > 0 { // OBP1 - 0xFF49
						colorIdx = getPaletteBits(f.ppu.gb.bus.read(0xFF49), colorIdx)
					} else { // OBP0 - 0xFF48
						colorIdx = getPaletteBits(f.ppu.gb.bus.read(0xFF48), colorIdx)
					}
				}

				// we decide the obj vs bg/win priority here
				// if obj wins, we push the colorIdx
				// if bg/win wins, we push transparent
				if getBit(obj.flags, 7) > 0 {
					objContest = true
				}
				objColorIdx = colorIdx
				objX = obj.x
			}
		}

		lo := getBit(f.targetLo, uint8(i))
		hi := getBit(f.targetHi, uint8(i))
		colorIdx := hi<<1 | lo

		if !f.ppu.gb.lcd.LCDC_bgwindowenable() {
			colorIdx = 0
		}

		colorIdx = getPaletteBits(f.ppu.gb.bus.read(0xFF47), colorIdx)

		f.ppu.fifo.bgQueue.enqueue(&fifoPixel{
			color: colorIdx,
		})

		if objContest {
			if colorIdx > 0 {
				objColorIdx = 0
			}
		}

		f.ppu.fifo.objQueue.enqueue(&fifoPixel{
			color: objColorIdx,
		})

		f.fetcherX += 1
	}

	f.state = fetcherStateReadTile
}
