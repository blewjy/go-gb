package gb

type fifoPixel struct {
	color   uint8
	palette uint8
}

type FIFO struct {
	ppu *ppu

	bgQueue  *queue
	objQueue *queue

	pxPushed uint8
}

func newFIFO(ppu *ppu) *FIFO {
	return &FIFO{
		ppu:      ppu,
		bgQueue:  newQueue(),
		objQueue: newQueue(),
	}
}

func (f *FIFO) pushOut() *fifoPixel {
	if f.bgQueue.len() != f.objQueue.len() {
		panic("nooo")
	}

	if f.bgQueue.len() <= 8 {
		return nil
	}

	bgPx, ok := f.bgQueue.dequeue()
	if !ok {
		panic("shit")
	}

	objPx, ok := f.objQueue.dequeue()
	if !ok {
		panic("shit")
	}

	f.pxPushed += 1

	// some assertion
	if f.pxPushed > 160 {
		panic("sum ting wong")
	}

	if objPx.color != 0 {
		return &fifoPixel{
			color:   objPx.color,
			palette: objPx.palette,
		}
	}
	return &fifoPixel{
		color:   bgPx.color,
		palette: bgPx.palette,
	}
}

func (f *FIFO) clear() {
	f.bgQueue.clear()
	f.objQueue.clear()
	f.pxPushed = 0
}

type node struct {
	val  *fifoPixel
	next *node
}

type queue struct {
	start *node
	end   *node
}

func newQueue() *queue {
	return &queue{}
}

func (q *queue) enqueue(val *fifoPixel) {
	newNode := &node{val: val, next: nil}
	if q.end == nil {
		q.start = newNode
		q.end = newNode
	} else {
		q.end.next = newNode
		q.end = newNode
	}
}

func (q *queue) dequeue() (*fifoPixel, bool) {
	if q.start == nil {
		return nil, false
	}
	item := q.start.val
	q.start = q.start.next
	if q.start == nil {
		q.end = nil
	}
	return item, true
}

func (q *queue) empty() bool {
	return q.start == nil
}

func (q *queue) clear() {
	q.start = nil
	q.end = nil
}

func (q *queue) len() int {
	if q.start == nil {
		return 0
	}
	count := 1
	nxt := q.start.next
	for nxt != nil {
		count += 1
		nxt = nxt.next
	}
	return count
}
