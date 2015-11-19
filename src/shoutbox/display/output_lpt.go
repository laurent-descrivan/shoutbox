package display

import (
	"os"
	"time"
)

const LPT_TICK = 0 * time.Nanosecond

type LptOutput struct {
	lastBuffer *Buffer
	buffer     *Buffer
	out        *os.File
}

func NewLptOutput(buffer *Buffer) *LptOutput {
	out, err := os.OpenFile("/dev/port", os.O_RDWR, 0777)
	if err != nil {
		panic(err.Error())
	}

	lastBuffer := NewBuffer()
	lastBuffer.CopyFrom(buffer)
	return &LptOutput{
		buffer:     buffer,
		lastBuffer: lastBuffer,
		out:        out,
	}
}

func (this *LptOutput) Buffer() *Buffer {
	return this.buffer
}

func (this *LptOutput) Flush() {
	newx := 0
	for oldx := 0; oldx < this.buffer.Width; oldx++ {
		if sameRows(this.buffer.Pixels[newx], this.lastBuffer.Pixels[oldx]) {
			newx += 1
		} else {
			newx = 0
		}
	}

	for x := newx; x < this.buffer.Width; x++ {
		var col byte
		height := uint(this.buffer.Height)
		for y := uint(0); y < height; y++ {
			if this.buffer.Pixels[x][y] {
				col |= 1 << (height - y - 1)
			}
		}
		this.putline(col)
	}
	this.lastBuffer.CopyFrom(this.buffer)
}

func (this *LptOutput) Clear() {
	for i := 0; i < this.buffer.Width; i++ {
		this.putline(0)
	}
}

func (this *LptOutput) putline(pixels byte) {
	this.out.WriteAt([]byte{(pixels & 127) ^ 255}, 888)
	time.Sleep(LPT_TICK)
	this.out.WriteAt([]byte{(pixels | 128) ^ 255}, 888)
	time.Sleep(LPT_TICK)
	this.out.WriteAt([]byte{(pixels & 127) ^ 255}, 888)
	time.Sleep(LPT_TICK)
}

func sameRows(rowa []bool, rowb []bool) bool {
	if len(rowa) != len(rowb) {
		return false
	}
	for i := range rowa {
		if rowa[i] != rowb[i] {
			return false
		}
	}
	return true
}
