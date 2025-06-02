package examples

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ProjectLighthouseCAU/lighthouse-go/colors"
	"github.com/ProjectLighthouseCAU/lighthouse-go/lighthouse"
)

func DisplayAPI(user, token, url string, fps int) {
	// Create a new display
	d, err := lighthouse.NewDisplay(user, token, url)
	if err != nil {
		log.Println(err)
		return
	}

	// Start the stream
	// stream, err := d.StartStream()
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// go func() {
	// 	for {
	// 		img, ok := <-stream
	// 		if !ok {
	// 			log.Println("stream closed!")
	// 			return
	// 		}
	// 		// log.Println("first pixel from stream: ", img[0:3])
	// 		_ = img
	// 	}
	// }()
	// Stop stream again after 3 seconds
	// time.AfterFunc(3*time.Second, func() {
	// 	err := d.StopStream()
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// })

	// Create a ticker to run code at fps times per second
	ticker := time.NewTicker(time.Second / time.Duration(fps))
	defer ticker.Stop()
	// Catch interrupts (e.g. Ctrl+C)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// choose animation:
	// update := rainbow(255)
	// update := rbgBlink(fps)
	// update := white(0)
	// update := color([3]byte{0, 255, 0})
	// update := rampUp(true, true, true)
	update := scanLine()
	for {
		select {
		case <-interrupt: // close display on interrupt
			d.Close()
			return
		case <-ticker.C: // on tick, update image and send
			img := update()
			err := d.SendImage(img)
			if err != nil {
				log.Println(err)
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
			log.Println(str)
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

func rampUp(r bool, g bool, b bool) func() []byte {
	image := make([]byte, 28*14*3)
	return func() []byte {
		for i := 0; i < len(image); i++ {
			switch i % 3 {
			case 0:
				if r {
					image[i] += 1
				}
			case 1:
				if g {
					image[i] += 1
				}
			case 2:
				if b {
					image[i] += 1
				}
			}
		}
		log.Println("Color: ", image[0])
		return image
	}
}

func scanLine() func() []byte {
	x := 0
	return func() []byte {
		image := make([]byte, 28*14*3)
		for y := 0; y < 14; y++ {
			image[(y*28+x)*3] = 255
			image[(y*28+x)*3+1] = 255
			image[(y*28+x)*3+2] = 255
		}
		x = (x + 1) % 28
		return image
	}
}
