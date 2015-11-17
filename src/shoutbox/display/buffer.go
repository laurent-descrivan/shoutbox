package display

type Buffer struct {
	Pixels [][]bool
	Width  int
	Height int
}

func NewBuffer() *Buffer {
	w := 96
	h := 7
	px := make([][]bool, w)
	for x := 0; x < w; x++ {
		px[x] = make([]bool, h)
	}
	return &Buffer{Pixels: px, Width: w, Height: h}
}

func (this *Buffer) CopyFrom(src *Buffer) {
	for x := 0; x < src.Width; x++ {
		for y := 0; y < src.Height; y++ {
			this.Pixels[x][y] = src.Pixels[x][y]
		}
	}
}
