package internal

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTTPSDPServer starts a HTTP Server that consumes SDPs
func HTTPSDPServer(port string) chan string {
	sdpChan := make(chan string)
	http.HandleFunc("/sdp", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		w.Header().Set("Access-Control-Allow-Origin", "localhost:8080")
		w.Write([]byte("done"))
		sdpChan <- string(body)
	})

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
		if err != nil {
			panic(err)
		}
	}()

	return sdpChan
}
