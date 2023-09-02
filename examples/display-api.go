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
	// stream, err := d.StartStream()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// go func() {
	// 	for {
	// 		img, ok := <-stream
	// 		if !ok {
	// 			fmt.Println("stream closed!")
	// 			return
	// 		}
	// 		// fmt.Println("first pixel from stream: ", img[0:3])
	// 		_ = img
	// 	}
	// }()
	// Stop stream again after 3 seconds
	// time.AfterFunc(3*time.Second, func() {
	// 	err := d.StopStream()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// })

	// Create a ticker to run code 60 times per second
	const fps int = 60
	ticker := time.NewTicker(time.Second / time.Duration(fps))
	defer ticker.Stop()
	// Catch interrupts (e.g. Ctrl+C)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// choose animation:
	update := rainbow(255)
	// update := rbgBlink(fps)
	// update := white(255)
	// update := color([3]byte{255, 0, 255})
	for {
		select {
		case <-interrupt: // close display on interrupt
			d.Close()
			return
		case <-ticker.C: // on tick, update image and send
			img := update()
			err := d.SendImage(img)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// brightness 0-255
func rainbow(brightness byte) func() []byte {
	// Create an image
	img := make([]byte, 28*14*3)

	// Some animation code (e.g. hsl color circle)
	var theta float64 = 0
	var xCurr int = 0
	var yCurr int = 0
	return func() []byte {
		for x := 0; x < 28; x++ {
			for y := 0; y < 14; y++ {
				_r, _g, _b := colors.HslToRgb(theta, 1, 0.5)
				r := byte(_r * float64(brightness))
				g := byte(_g * float64(brightness))
				b := byte(_b * float64(brightness))
				i := 3 * (x + y*28)
				img[i] = r
				img[i+1] = g
				img[i+2] = b
			}
		}
		theta += 0.0005
		xCurr = (xCurr + 1) % 28
		yCurr = (yCurr + 1) % 14
		if theta >= 1 {
			theta = 0
		}
		return img
	}
}

func rbgBlink(threshold int) func() []byte {
	red := make([]byte, 28*14*3)
	green := make([]byte, 28*14*3)
	blue := make([]byte, 28*14*3)
	images := [][]byte{red, green, blue}
	idx := 0
	counter := 0
	return func() []byte {
		for i := 0; i < 28*14*3; i += 3 {
			red[i] = 255
			green[i+1] = 255
			blue[i+2] = 255
		}

		counter++
		if counter == threshold {
			idx = (idx + 1) % len(images)
			counter = 0
			var str string
			if idx == 0 {
				str = "RED"
			} else if idx == 1 {
				str = "GREEN"
			} else if idx == 2 {
				str = "BLUE"
			}
			fmt.Println(str)
		}
		return images[idx]
	}
}

func white(brightness byte) func() []byte {
	return color([3]byte{brightness, brightness, brightness})
}

func color(color [3]byte) func() []byte {
	image := make([]byte, 28*14*3)
	for i := 0; i < 28*14*3; i += 3 {
		image[i] = color[0]
		image[i+1] = color[1]
		image[i+2] = color[2]
	}
	return func() []byte {
		return image
	}
}
