package examples

import (
	"log"
	"sync"

	"github.com/ProjectLighthouseCAU/lighthouse-go/lighthouse"
)

func ClientAPI(user, token, url string) {
	// Create a new client
	client, err := lighthouse.NewClient(url)
	if err != nil {
		log.Println(err)
		return
	}
	// Create a request
	request := lighthouse.NewRequest().Reid(0).Auth(user, token).Verb("GET").Path("user", user, "model")

	// Waitgroup to wait for all responses
	var wg sync.WaitGroup
	reids := make(map[int8]struct{})

	// Goroutine to handle the responses
	go func() {
		for {
			resp, err := client.Receive()
			if err != nil {
				log.Println(err)
				break
			}
			delete(reids, resp.REID.(int8))
			wg.Done()
			log.Printf("Response: %+v\n", resp)
		}
	}()

	// Send some requests
	var i int8
	var n int8 = 10
	for i = 0; i < n; i++ {
		reids[i] = struct{}{}
	}
	for i = 0; i < n; i++ {
		request.Reid(i)
		log.Printf("Request: %+v\n", request)
		wg.Add(1)
		client.Send(request)
	}
	// Wait for all responses
	wg.Wait()
	client.Close()

	for i := range reids {
		log.Printf("No response for %d\n", i)
	}
}
