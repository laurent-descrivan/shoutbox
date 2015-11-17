package display

import (
	"os"
	"time"
)

const LPT_TICK = 1000 * time.Nanosecond

type LptOutput struct {
	buffer *Buffer
	out    *os.File
}

func NewLptOutput(buffer *Buffer) *LptOutput {
	out, err := os.OpenFile("/dev/port", os.O_RDWR, 0777)
	if err != nil {
		panic(err.Error())
	}
	return &LptOutput{
		buffer: buffer,
		out:    out,
	}
}

func (this *LptOutput) Buffer() *Buffer {
	return this.buffer
}

func (this *LptOutput) Flush() {
	for x := 0; x < this.buffer.Width; x++ {
		var col byte
		height := uint(this.buffer.Height)
		for y := uint(0); y < height; y++ {
			if this.buffer.Pixels[x][y] {
				col |= 1 << (height - y - 1)
			}
		}
		this.putline(col)
	}
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
