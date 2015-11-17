package main

import (
	"image/gif"
	"os"
	"shoutbox/animation"
	"shoutbox/display"
)

func main() {
	buffer := display.NewBuffer()
	output := display.NewLptOutput(buffer) // display.NewTerminalOutput(buffer)

	animator := animation.NewGifAnimator(output)

	f, err := os.Open("data/2.gif")
	if err != nil {
		panic(err.Error())
	}
	g, err := gif.DecodeAll(f)
	if err != nil {
		panic(err.Error())
	}

	animator.SetGif(g)
	animator.Start()
	animator.Wait()
}
