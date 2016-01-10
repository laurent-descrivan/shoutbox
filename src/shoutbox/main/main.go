package main

import (
	"fmt"
	// "image/gif"
	"net"
	// "os"
	"shoutbox/animation"
	"shoutbox/display"
	"time"
)

func main() {
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
	time.Sleep(15 * time.Second)

	animator.SetText("-= Electrolab =-")
	animator.Wait()
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
