package main

import (
	"flag"

	"github.com/ProjectLighthouseCAU/lighthouse-go/examples"
)

var (
	user  string
	token string
	url   string
)

func init() {
	flag.StringVar(&user, "user", "Testuser", "username")
	flag.StringVar(&token, "token", "API-TOK_TEST", "API token")
	flag.StringVar(&url, "url", "ws://localhost:3000/websocket", "websocket url (ws:// or wss://)")
}

func main() {
	flag.Parse()
	examples.DisplayAPI(user, token, url)
	// examples.ClientAPI(user, token, url)
}
