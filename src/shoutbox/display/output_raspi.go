package display

import (
	"github.com/stianeikeland/go-rpio"
	"time"
)

const RASPI_TICK = 0 * time.Nanosecond

type RaspiOutput struct {
	lastBuffer *Buffer
	buffer     *Buffer
	clock      rpio.Pin
	out        [7]rpio.Pin
}

func NewRaspiOutput(buffer *Buffer) *RaspiOutput {
	err := rpio.Open()
	if err != nil {
		panic(err.Error())
	}

	lastBuffer := NewBuffer()
	lastBuffer.CopyFrom(buffer)
	return &RaspiOutput{
		buffer:     buffer,
		lastBuffer: lastBuffer,
		clock:      rpio.Pin(15),
		out: [7]rpio.Pin{
			rpio.Pin(0),
			rpio.Pin(1),
			rpio.Pin(4),
			rpio.Pin(17),
			rpio.Pin(21),
			rpio.Pin(22),
			rpio.Pin(14),
		},
	}
}

func (this *RaspiOutput) Buffer() *Buffer {
	return this.buffer
}

func (this *RaspiOutput) Flush() {
	// Look at the differences between last frame and current frame,
	// and compute the number of shifted columns, in order to push
	// only the minimum amount of pixel columns.
	var shiftedColumnsCount = 0
	for shiftedColumnsCount = 0; shiftedColumnsCount < this.buffer.Width; shiftedColumnsCount++ {
		if sameRows(this.lastBuffer.Pixels[shiftedColumnsCount:], this.buffer.Pixels[:this.buffer.Width-shiftedColumnsCount]) {
			break
		}
	}

	// If the last frame is the same, do nothing
	if shiftedColumnsCount == 0 {
		return
	}

	// Push the shifted columns, up to the entire width of the buffer
	for x := this.buffer.Width - shiftedColumnsCount; x < this.buffer.Width; x++ {
		var col byte
		height := uint(this.buffer.Height)
		for y := uint(0); y < height; y++ {
			if this.buffer.Pixels[x][y] {
				col |= 1 << (height - y - 1)
			}
		}
		this.putline(col)
	}

	// Backup the frame for later comparison
	this.lastBuffer.CopyFrom(this.buffer)
}

func (this *RaspiOutput) Clear() {
	for i := 0; i < this.buffer.Width; i++ {
		this.putline(0)
	}
}

func (this *RaspiOutput) putline(pixels byte) {
	for i := uint(0); i < 7; i++ {
		setPin(this.out[i], (pixels|(1<<i)) != 0)
	}
	setPin(this.clock, true)
	time.Sleep(10 * time.Microsecond)
	setPin(this.clock, false)
	time.Sleep(10 * time.Microsecond)
}

func setPin(pin rpio.Pin, state bool) {
	if state {
		pin.High()
	} else {
		pin.Low()
	}
}
