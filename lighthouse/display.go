package lighthouse

import (
	"errors"
	"fmt"
)

const (
	channelSize = 100
)

// Display is a simple API wrapper for easy animation or game development
// Note that the methods are not thread-safe!
type Display struct {
	client  *Client
	request *Request
	stream  chan []byte
}

// Create a new display
func NewDisplay(user string, token string, url string) (*Display, error) {
	c, err := NewClient(url)
	if err != nil {
		return nil, err
	}
	d := &Display{
		client:  c,
		request: NewRequest().Reid(0).Auth(user, token).Verb("POST").Path("user", user, "model"),
		stream:  make(chan []byte, channelSize),
	}
	go d.responseHandler()
	return d, nil
}

// Closes the display and the underlying client
// You cannot open the display again but instead create a new one
func (d *Display) Close() {
	close(d.stream)
	d.client.Close()
}

// Sends an image
func (d *Display) SendImage(img []byte) error {
	if len(img) != 28*14*3 {
		return errors.New("SendImage: Image ([]byte) has wrong size (must be 28*14*3)")
	}
	return d.client.Send(d.request.Reid(0).Verb("POST").Payl(img))
}

// Starts the stream and returns a read-only channel containing the images from the stream
func (d *Display) StartStream() (<-chan []byte, error) {
	if err := d.client.Send(d.request.Reid(1).Verb("STREAM").Payl(nil)); err != nil {
		return nil, err
	}
	return d.stream, nil
}

// Stops the stream
func (d *Display) StopStream() error {
	return d.client.Send(d.request.Reid(1).Verb("STOP").Payl(nil))
}

// Goroutine for handling the responses
func (d *Display) responseHandler() {
	for {
		resp, err := d.client.Receive()
		if err != nil {
			fmt.Println(err)
			d.Close()
			return
		}
		reid, ok := resp.REID.(int8)
		if !ok {
			continue
		}
		switch reid {
		case 0: // PUT/POST response
			if resp.RNUM >= 400 { // print only errors
				fmt.Printf("%+v\n", resp)
			}
		case 1: // STREAM response
			if resp.RNUM >= 400 { // print only errors
				fmt.Printf("%+v\n", resp)
			}
			// forward to image stream
			payl, ok := resp.PAYL.([]byte)
			if !ok || len(payl) != 28*14*3 {
				continue
			}
			select {
			case d.stream <- payl:
			default:
				fmt.Println("Stream channel full, skipping")
			}
		}
	}
}
