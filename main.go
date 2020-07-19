package main

import (
	"github.com/riyadennis/chatterbox/internal"
	"github.com/riyadennis/chatterbox/internal/handler"
)

func main() {
	sdpChan := internal.HTTPSDPServer(handler.Port)
	err := handler.SDPRequest(<-sdpChan)
	if err != nil {
		panic(err)
	}
}
