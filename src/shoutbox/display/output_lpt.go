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
	var count = 0
	for count := 0; count < this.buffer.Width; count++ {
		if sameRows(this.lastBuffer.Pixels[count:], this.buffer.Pixels[:this.buffer.Width-count]) {
			break
		}
	}

	for x := count; x < this.buffer.Width; x++ {
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

func sameRows(rowsa, rowsb [][]bool) bool {
	for i := range rowsa {
		rowa := rowsa[i]
		rowb := rowsb[i]
		for j := range rowa {
			if rowa[j] != rowb[j] {
				return false
			}
		}
	}
	return false
}
