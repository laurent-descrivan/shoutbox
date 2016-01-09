package animation

import (
	"image/color"
	"image/gif"
	"shoutbox/display"
	"sync"
	"time"
)

type GifAnimator struct {
	output display.Output
	g      *gif.GIF

	isRunning     bool
	orderStop     chan bool
	waitForFinish sync.WaitGroup
}

func NewGifAnimator(output display.Output) *GifAnimator {
	return &GifAnimator{
		output: output,
	}
}

func (this *GifAnimator) SetGif(g *gif.GIF) {
	wasRunning := this.isRunning
	if wasRunning {
		this.Stop()
	}
	this.g = g
	if wasRunning {
		this.Start()
	}
}

func (this *GifAnimator) Start() {
	if !this.isRunning {
		this.orderStop = make(chan bool)
		this.isRunning = true
		this.waitForFinish.Add(1)
		go this.loop()
	}
}

func (this *GifAnimator) Stop() {
	if this.isRunning {
		close(this.orderStop)
		this.waitForFinish.Wait()
	}
}

func (this *GifAnimator) Wait() {
	this.waitForFinish.Wait()
}

func (this *GifAnimator) loop() {
	defer this.waitForFinish.Done()
	if this.g == nil {
		return
	}

	buffer := this.output.Buffer()
	restorePreviousBuffer := display.NewBuffer()

	for loopCount := 0; this.g.LoopCount == 0 || loopCount < this.g.LoopCount; loopCount++ {
		for idx := 0; idx < len(this.g.Image); idx++ {
			img := this.g.Image[idx]
			// bgColor := img.Palette[this.g.BackgroundIndex].RGBA()
			// imgDisposal := this.g.Disposal[i]

			for x := 0; x < buffer.Width; x++ {
				for y := 0; y < buffer.Height; y++ {
					c := img.At(x, y)
					_, _, _, alpha := img.At(x, y).RGBA()
					if alpha >= 0x8000 {
						buffer.Pixels[x][y] = isLit(c)
					} else {
						// See http://www.webreference.com/content/studio/disposal.html
						switch this.g.Disposal[idx] {
						case gif.DisposalBackground: // Restore to Background if asked for
							buffer.Pixels[x][y] = isLit(img.Palette[this.g.BackgroundIndex])
						case gif.DisposalPrevious: // Restore to previous if asked for
							buffer.Pixels[x][y] = restorePreviousBuffer.Pixels[x][y]
						case gif.DisposalNone:
						default:
							buffer.Pixels[x][y] = false
						}
					}
				}
			}

			// Backup the buffer in case DisposalPrevious is used later
			switch this.g.Disposal[idx] {
			case 0: // Unspecified
				fallthrough
			case gif.DisposalNone: // Do Not Dispose
				restorePreviousBuffer.CopyFrom(buffer)
			}

			this.output.Flush()

			delay := 10 * time.Millisecond * time.Duration(this.g.Delay[idx])
			if delay <= 0 {
				delay = 100000 * time.Hour
			}

			select {
			case <-this.orderStop:
				return
			case <-time.After(delay):
			}
		}
	}
}

func isLit(c color.Color) bool {
	r, g, b, a := c.RGBA()
	if a < 0x8000 {
		return true
	} else {
		return ((0.2126*float32(r) + 0.7152*float32(g) + 0.0722*float32(b)) / 0xffff) < 0.5
	}
}
