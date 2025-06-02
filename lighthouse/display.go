package lighthouse

import (
	"errors"
	"log"
	"strings"
	"sync"
)

const (
	defaultChannelSize = 1
)

// Display is a simple API wrapper for easy animation or game development
// Note that the methods are not thread-safe!
type Display struct {
	client      *Client
	User        string
	Token       string
	channelSize int
	streams     map[string]chan any
	streamsLock sync.Mutex
}

// Create a new display
func NewDisplay(user string, token string, url string) (*Display, error) {
	return NewDisplayWithChannelSize(user, token, url, defaultChannelSize)
}

func NewDisplayWithChannelSize(user string, token string, url string, channelSize int) (*Display, error) {
	c, err := NewClient(url)
	if err != nil {
		return nil, err
	}
	d := &Display{
		client:      c,
		User:        user,
		Token:       token,
		channelSize: channelSize,
		streams:     make(map[string]chan any),
	}
	go d.responseHandler()
	return d, nil
}

// Closes the display and the underlying client
// You cannot open the display again but instead create a new one
func (d *Display) Close() {
	d.client.Close()
}

// Returns a pointer to the underlying client
func (d *Display) GetClient() *Client {
	return d.client
}

// Sends an image
func (d *Display) SendImage(img []byte) error {
	if len(img) != 28*14*3 {
		return errors.New("SendImage: Image ([]byte) has wrong size (must be 28*14*3)")
	}
	return d.client.Send(NewRequest().Auth(d.User, d.Token).Path("user", d.User, "model").Reid(0).Verb("PUT").Payl(img))
}

// Starts the stream and returns a read-only channel containing the images from the stream
func (d *Display) Stream(path []string) (<-chan any, error) {
	pathStr := strings.Join(path, "/")
	if err := d.client.Send(NewRequest().Auth(d.User, d.Token).Path(path...).Reid(pathStr).Verb("STREAM").Payl(nil)); err != nil {
		return nil, err
	}
	c := make(chan any, d.channelSize)
	d.streamsLock.Lock()
	d.streams[pathStr] = c
	d.streamsLock.Unlock()
	return c, nil
}

// Stops the stream
func (d *Display) StopStream(path []string) error {
	pathStr := strings.Join(path, "/")
	err := d.client.Send(NewRequest().Auth(d.User, d.Token).Path(path...).Reid(pathStr).Verb("STOP").Payl(nil))

	d.streamsLock.Lock()
	stream, ok := d.streams[pathStr]
	if ok {
		delete(d.streams, pathStr)
		close(stream)
	}
	d.streamsLock.Unlock()

	return err
}

// Goroutine for handling the responses
func (d *Display) responseHandler() {
	for {
		resp, err := d.client.Receive()
		if err != nil {
			log.Println(err)
			for _, stream := range d.streams {
				close(stream)
			}
			return
		}
		// interpret REID as path string
		key, ok := resp.REID.(string)
		if !ok {
			if resp.RNUM >= 400 { // log only errors
				log.Printf("%+v\n", resp)
			}
			continue
		}
		// get stream
		d.streamsLock.Lock()
		stream, ok := d.streams[key]
		if !ok {
			d.streamsLock.Unlock()
			continue
		}
		// close and delete stream on error and log error
		if resp.RNUM >= 400 {
			log.Printf("%+v\n", resp)
			delete(d.streams, key)
			close(stream)
			d.streamsLock.Unlock()
			continue
		}
		// send payload to stream
		stream <- resp.PAYL
		d.streamsLock.Unlock()
	}
}
