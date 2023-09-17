package gb

type oamEntry struct {
	y        uint8
	x        uint8
	tileIdx  uint8
	flags    uint8 // todo: can expand to individual bits
	size     uint8
	fetchedX uint8
}

type ppu struct {
	gb *Gameboy

	oamLineEntries []*oamEntry
	oamScanCount   uint64

	fetcher *fetcher
	fifo    *FIFO

	dots uint64

	framePx      []*fifoPixel
	framePxReady []*fifoPixel
}

func newPpu(gb *Gameboy) *ppu {
	p := &ppu{
		gb:   gb,
		dots: 0,
	}
	p.fetcher = newFetcher(p)
	p.fifo = newFIFO(p)
	return p
}

func (p *ppu) tick() {
	switch p.gb.lcd.get_STAT_modeflag() {
	case STAT_mode_searchoam:
		p.tick_oam()
	case STAT_mode_transferring:
		p.tick_transferring()
	case STAT_mode_hblank:
		p.tick_hblank()
	case STAT_mode_vblank:
		p.tick_vblank()
	}
	p.dots++
}

func (p *ppu) tick_oam() {
	if p.dots%2 == 0 {
		if len(p.oamLineEntries) < 10 && p.gb.lcd.LCDC_objenable() { // For CGB, no need to check objenable()
			addr := uint16(0xFE00) + uint16(p.oamScanCount)*4
			oam_entry := &oamEntry{
				y:       p.gb.bus.read(addr),
				x:       p.gb.bus.read(addr + 1),
				tileIdx: p.gb.bus.read(addr + 2),
				flags:   p.gb.bus.read(addr + 3),
				size:    0,
			}

			// determine fetchedX based on whether flipped or not
			oam_entry.fetchedX = 7
			if getBit(oam_entry.flags, 5) > 0 {
				oam_entry.fetchedX = 0
			}

			y_end := oam_entry.y + 7
			if p.gb.lcd.LCDC_objsize() > 0 {
				oam_entry.size = 1
				y_end = oam_entry.y + 15

				y_flipped := getBit(oam_entry.flags, 6) > 0

				if !y_flipped {
					if p.gb.lcd.LY() >= oam_entry.y-8 {
						oam_entry.tileIdx |= 0x01
					} else {
						oam_entry.tileIdx &= 0xFE
					}
				} else {
					if p.gb.lcd.LY() >= oam_entry.y-8 {
						oam_entry.tileIdx &= 0xFE
					} else {
						oam_entry.tileIdx |= 0x01
					}
				}
			}

			if p.gb.lcd.LY() >= oam_entry.y-16 && p.gb.lcd.LY() <= y_end-16 {
				p.oamLineEntries = append(p.oamLineEntries, oam_entry)
			}
		}
		p.oamScanCount++
	}

	// some assertions
	if p.oamScanCount > 40 || len(p.oamLineEntries) > 10 {
		panic(p.oamScanCount)
	}

	if p.dots >= 79 {
		p.gb.lcd.set_STAT_modeflag(STAT_mode_transferring)
		p.fifo.clear()
		p.fetcher.state = fetcherStateReadTile
		p.fetcher.fetcherX = 0
		p.fetcher.prevInWindow = false
	}
}

func (p *ppu) tick_transferring() {
	px := p.fifo.pushOut()
	if px != nil {
		p.framePx = append(p.framePx, px)
	}

	if p.dots%2 == 0 {
		p.fetcher.clock()
	}

	// once pushed 160 px (width of lcd), can move on to hblank
	if p.fifo.pxPushed >= 160 {
		p.gb.lcd.set_STAT_modeflag(STAT_mode_hblank)

		// going from xfer -> hblank
		if p.gb.lcd.STAT_hblankinterruptsource() {
			p.gb.cpu.reqInterrupt(cpuInterruptLcdStat)
		}

	}
}

func (p *ppu) tick_hblank() {
	if p.dots >= 455 {
		p.gb.lcd.incLY()
		if p.fetcher.windowDidDisplay {
			p.fetcher.windowY += 1
		}
		p.fetcher.windowDidDisplay = false

		p.gb.lcd.set_STAT_lycflag(p.gb.lcd.LY() == p.gb.lcd.LYC())
		if p.gb.lcd.get_STAT_lycflag() && p.gb.lcd.STAT_lycinterruptsource() {
			p.gb.cpu.reqInterrupt(cpuInterruptLcdStat)
		}

		if p.gb.lcd.LY() >= 144 {
			p.gb.lcd.set_STAT_modeflag(STAT_mode_vblank)
			p.gb.cpu.reqInterrupt(cpuInterruptVBlank)

			// going from hblank -> vblank
			if p.gb.lcd.STAT_vblankinterruptsource() {
				p.gb.cpu.reqInterrupt(cpuInterruptLcdStat)
			}

		} else {
			p.gb.lcd.set_STAT_modeflag(STAT_mode_searchoam)
			p.oamScanCount = 0
			p.oamLineEntries = nil
			// going from hblank -> oam search
			if p.gb.lcd.STAT_oaminterruptsource() {
				p.gb.cpu.reqInterrupt(cpuInterruptLcdStat)
			}
		}

		p.dots = 0
	}
}

func (p *ppu) tick_vblank() {
	if p.dots >= 455 {
		p.gb.lcd.incLY()
		if p.fetcher.windowDidDisplay {
			p.fetcher.windowY += 1
		}
		p.fetcher.windowDidDisplay = false

		p.gb.lcd.set_STAT_lycflag(p.gb.lcd.LY() == p.gb.lcd.LYC())
		if p.gb.lcd.get_STAT_lycflag() && p.gb.lcd.STAT_lycinterruptsource() {
			p.gb.cpu.reqInterrupt(cpuInterruptLcdStat)
		}

		if p.gb.lcd.LY() >= 153 {
			p.gb.lcd.set_STAT_modeflag(STAT_mode_searchoam)
			p.oamScanCount = 0
			p.oamLineEntries = nil

			p.framePxReady = make([]*fifoPixel, len(p.framePx))
			copy(p.framePxReady, p.framePx)
			p.framePx = nil

			p.gb.lcd.resetLY()
			p.fetcher.windowY = 0

			// going from hblank -> oam search
			if p.gb.lcd.STAT_oaminterruptsource() {
				p.gb.cpu.reqInterrupt(cpuInterruptLcdStat)
			}
		}

		p.dots = 0
	}
}
