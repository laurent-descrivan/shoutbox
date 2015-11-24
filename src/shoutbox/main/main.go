package main

import (
	"image/gif"
	"os"
	"shoutbox/animation"
	"shoutbox/display"
)

func main() {
	buffer := display.NewBuffer()
	output := display.NewLptOutput(buffer)
	// output := display.NewTerminalOutput(buffer)

	animator := animation.NewGifAnimator(output)
	// animator := animation.NewTextAnimator(output)

	f, err := os.Open("data/pacman.gif")
	if err != nil {
		panic(err.Error())
	}
	g, err := gif.DecodeAll(f)
	if err != nil {
		panic(err.Error())
	}

	animator.SetGif(g)

	// animator.SetText("-= Electrolab =-")
	animator.Start()
	animator.Wait()
}
