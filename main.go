package main

import (
	"github.com/riyadennis/chatterbox/internal"
	"github.com/riyadennis/chatterbox/internal/handler"
	"github.com/sirupsen/logrus"
)

func main() {
	sdpChan, errChan := internal.HTTPSDPServer(handler.Port)
	for {
		select {
		case err := <-errChan:
			if err != nil {
				logrus.Fatalf("error from request :: %v", err)
			}
		case sdp := <-sdpChan:
			handler.SDPRequest(sdp, errChan)
		}
	}
}
