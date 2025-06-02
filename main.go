package main

import (
	"flag"

	"github.com/ProjectLighthouseCAU/lighthouse-go/examples"
)

var (
	user  string
	token string
	url   string
	fps   int
)

func init() {
	flag.StringVar(&user, "user", "", "username")
	flag.StringVar(&token, "token", "", "API token")
	flag.StringVar(&url, "url", "wss://lighthouse.uni-kiel.de/websocket", "websocket url (ws:// or wss://)")
	flag.IntVar(&fps, "fps", 60, "fps")
}

func main() {
	flag.Parse()
	examples.DisplayAPI(user, token, url, fps)
	// examples.ClientAPI(user, token, url)
}
