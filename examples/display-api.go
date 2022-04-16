package examples

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/ProjectLighthouseCAU/lighthouse-go/colors"
	"github.com/ProjectLighthouseCAU/lighthouse-go/lighthouse"
)

func DisplayAPI(user, token, url string) {
	// Create a new display
	d, err := lighthouse.NewDisplay(user, token, url)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Start the stream
	stream, err := d.StartStream()
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		for {
			img, ok := <-stream
			if !ok {
				return
			}
			fmt.Println("first pixel from stream: ", img[0:3])
		}
	}()
	// Stop stream again after 3 seconds
	time.AfterFunc(3*time.Second, func() {
		d.StopStream()
	})

	// Create a ticker to run code 60 times per second
	ticker := time.NewTicker(time.Duration(1000/60) * time.Millisecond)
	defer ticker.Stop()
	// Catch interrupts (e.g. Ctrl+C)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Create an image
	img := make([]byte, 28*14*3)

	// Some animation code (e.g. hsl color circle)
	var theta float64 = 0
	var xCurr int = 0
	var yCurr int = 0
	for {
		select {
		case <-interrupt: // close display on interrupt
			d.Close()
			return
		case <-ticker.C: // on tick, update image and send
			for x := 0; x < 28; x++ {
				for y := 0; y < 14; y++ {
					_r, _g, _b := colors.HslToRgb(theta, 1, 0.5)
					r := byte(_r * 255)
					g := byte(_g * 255)
					b := byte(_b * 255)
					i := 3 * (x + y*28)
					img[i] = r
					img[i+1] = g
					img[i+2] = b
				}
			}
			err := d.SendImage(img)
			if err != nil {
				fmt.Println(err)
				return
			}
			theta += 0.005
			xCurr = (xCurr + 1) % 28
			yCurr = (yCurr + 1) % 14
			if theta >= 1 {
				theta = 0
			}
		}
	}
}
