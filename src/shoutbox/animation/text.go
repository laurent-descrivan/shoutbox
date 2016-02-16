package animation

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"image/draw"
	"io/ioutil"
	"shoutbox/display"
	"sync"
	"time"
)

const (
	TEXT_DELAY = 25 * time.Millisecond
	FONT_PATH  = "data/zig.ttf"
)

var ttf *truetype.Font

func init() {
	fontBytes, err := ioutil.ReadFile(FONT_PATH)
	if err != nil {
		panic(err.Error())
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err.Error())
	}
	ttf = f
}

type TextAnimator struct {
	output       display.Output
	text         string
	textImg      image.Image
	textImgWidth int

	isRunning     bool
	orderStop     chan bool
	waitForFinish sync.WaitGroup

	EndOfLine chan bool
}

func NewTextAnimator(output display.Output) *TextAnimator {
	return &TextAnimator{
		output:    output,
		EndOfLine: make(chan bool, 1),
	}
}

func (this *TextAnimator) SetText(text string) {
	this.text = text
	this.computeText()
}

func (this *TextAnimator) Start() {
	if !this.isRunning {
		this.orderStop = make(chan bool)
		this.isRunning = true
		this.waitForFinish.Add(1)
		go this.loop()
	}
}

func (this *TextAnimator) Stop() {
	if this.isRunning {
		close(this.orderStop)
		this.waitForFinish.Wait()
	}
}

func (this *TextAnimator) Wait() {
	this.waitForFinish.Wait()
}

func (this *TextAnimator) loop() {
	defer this.waitForFinish.Done()
	if len(this.text) == 0 {
		return
	}

	buffer := this.output.Buffer()

	for {
		for x := -buffer.Width; x < this.textImgWidth; x++ {
			this.copyImg(x, 0)
			this.output.Flush()
			select {
			case <-this.orderStop:
				return
			case <-time.After(TEXT_DELAY):
			}
		}
		select {
		case this.EndOfLine <- true:
		default:
		}
	}
}

func (this *TextAnimator) computeText() {
	buffer := this.output.Buffer()
	// Initialize the context.
	fg, bg := image.White, image.Black
	c := freetype.NewContext()
	c.SetDPI(100)
	c.SetFont(ttf)
	c.SetFontSize(float64(buffer.Height))
	c.SetSrc(fg)
	c.SetHinting(font.HintingFull)

	origin := freetype.Pt(0, buffer.Height)
	text := this.text + "  "

	var textWidth int
	if pt, err := c.DrawString(text, origin); err != nil {
		panic(err.Error())
	} else {
		textWidth = int(pt.X >> 6)
	}

	// Draw the text.

	rgba := image.NewRGBA(image.Rect(0, 0, textWidth, buffer.Height))
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)

	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

	if _, err := c.DrawString(text, origin); err != nil {
		panic(err.Error())
	}

	this.textImgWidth = textWidth
	this.textImg = rgba
}

func (this *TextAnimator) copyImg(x, y int) {
	buffer := this.output.Buffer()

	for x2 := 0; x2 < buffer.Width; x2++ {
		for y2 := 0; y2 < buffer.Height; y2++ {
			buffer.Pixels[x2][y2] = !isLit(this.textImg.At(x2+x, y2+y))
		}
	}
}
