package main

import (
	"fmt"
	"github.com/riyadennis/chatterbox/internal/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/sdp", handler.SDPRequest)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", handler.Port), nil))
}
