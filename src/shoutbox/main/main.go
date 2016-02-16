package main

import (
	"fmt"
	"net/http"
	// "image/gif"
	"net"
	"shoutbox/animation"
	"shoutbox/display"
	"shoutbox/server"
)

const LINES_FILE_PATH = `data/lines.txt`

func main() {
	server.ReloadText()

	mux := server.NewMux()
	go func() {
		err := http.ListenAndServe(":80", mux)
		if err != nil {
			panic(fmt.Sprintf("Cannot open web server: %s", err.Error()))
		}
	}()

	buffer := display.NewBuffer()
	// output := display.NewLptOutput(buffer)
	// output := display.NewTerminalOutput(buffer)
	output := display.NewRaspiOutput(buffer)

	// animator := animation.NewGifAnimator(output)
	animator := animation.NewTextAnimator(output)

	// f, err := os.Open("data/pacman.gif")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// g, err := gif.DecodeAll(f)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// _ = g
	// // animator.SetGif(g)

	ip := getNetworkIP()
	if ip == "" {
		ip = "<unknown ip>"
	}
	animator.SetText(fmt.Sprintf("http://%s/", ip))
	animator.Start()
	<-animator.EndOfLine
	<-animator.EndOfLine

	for {
		animator.SetText(server.GetNextLine())
		<-animator.EndOfLine
	}
	// animator.Wait()
}

func getNetworkIP() string {
	if addresses, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addresses {
			if ip, _, err := net.ParseCIDR(addr.String()); err == nil {
				if ip.To4() != nil && ip.IsGlobalUnicast() {
					return ip.String()
				}
			}
		}
	}

	return ""
}
