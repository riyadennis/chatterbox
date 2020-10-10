package internal

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTTPSDPServer starts a HTTP Server that consumes SDPs
func HTTPSDPServer(port string) (chan string, chan error) {
	sdpChan := make(chan string)
	errChan := make(chan error)
	http.HandleFunc("/sdp", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errChan <- err
		}
		w.Header().Set("Access-Control-Allow-Origin", "localhost:8080")
		w.Write([]byte("done"))
		sdpChan <- string(body)
	})

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
		if err != nil {
			errChan <- err
		}
	}()

	return sdpChan, errChan
}
