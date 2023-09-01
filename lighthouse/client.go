package lighthouse

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vmihailenco/msgpack"
)

// Client is the API wrapper for communicating with the lighthouse server
type Client struct {
	connection *websocket.Conn
}

// NewClient creates a new client and connects it to the lighthouse server at the given url
func NewClient(url string) (*Client, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	c := &Client{
		connection: conn,
	}
	fmt.Println("Connected")
	return c, nil
}

// Sends a request to the websocket connection
func (c *Client) Send(req *Request) error {
	encoded, err := msgpack.Marshal(req)
	if err != nil {
		return err
	}
	return c.connection.WriteMessage(websocket.BinaryMessage, encoded)
}

// Receives a response from the websocket connection if available or blocks until a response arrives
func (c *Client) Receive() (*Response, error) {
	_, encoded, err := c.connection.ReadMessage()
	if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		return nil, fmt.Errorf("websocket connection closed: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("error reading from websocket: %w", err)
	}
	resp := Response{}
	err = msgpack.Unmarshal(encoded, &resp)
	if err != nil {
		return nil, fmt.Errorf("unable to decode response: %w", err)
	}
	return &resp, nil
}

// Closes the connection gracefully
func (c *Client) Close() error {
	// Catch close message response
	closed := make(chan struct{})
	c.connection.SetPingHandler(nil)
	c.connection.SetPongHandler(nil)
	c.connection.SetCloseHandler(func(code int, text string) error {
		close(closed)
		return nil
	})
	c.connection.SetReadDeadline(time.Now().Add(time.Second))

	// Send close message
	err := c.connection.WriteControl(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "connection closed"),
		time.Now().Add(time.Second))

	if err != nil {
		c.connection.Close()
		return fmt.Errorf("could not send close message: %w", err)
	}
	go func() {
		for {
			_, _, err := c.connection.ReadMessage()
			if err != nil {
				c.connection.Close()
				return
			}
		}
	}()
	select {
	case <-closed:
		return c.connection.Close()
	case <-time.After(time.Second):
		return fmt.Errorf("did not receive close message: %w", c.connection.Close())
	}
}
