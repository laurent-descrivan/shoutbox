package display

import (
	"fmt"
	"strings"
)

type TerminalOutput struct {
	buffer       *Buffer
	firstDisplay bool
}

func NewTerminalOutput(buffer *Buffer) *TerminalOutput {
	return &TerminalOutput{
		buffer:       buffer,
		firstDisplay: true,
	}
}

func (this *TerminalOutput) Buffer() *Buffer {
	return this.buffer
}

func (this *TerminalOutput) Flush() {
	this.print(false)
}

func (this *TerminalOutput) Clear() {
	this.print(true)
}

func (this *TerminalOutput) print(clearing bool) {
	if this.firstDisplay {
		this.firstDisplay = false
	} else {
		fmt.Printf("\r\033[%dA", this.buffer.Height+2)
	}
	fmt.Printf("+%s+\n", strings.Repeat("-", this.buffer.Width))
	for y := 0; y < this.buffer.Height; y++ {
		fmt.Print("|")
		for x := 0; x < this.buffer.Width; x++ {
			if !this.buffer.Pixels[x][y] || clearing {
				fmt.Print(" ")
			} else {
				fmt.Print("*")
			}
		}
		fmt.Println("|")
	}
	fmt.Printf("+%s+\n", strings.Repeat("-", this.buffer.Width))
}
