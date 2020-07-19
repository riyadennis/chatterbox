package main

import (
	"github.com/riyadennis/chatterbox/internal"
	"github.com/riyadennis/chatterbox/internal/handler"
	"github.com/sirupsen/logrus"
)

func main() {
	errChan := make(chan error, 2)
	sdpChan := internal.HTTPSDPServer(handler.Port)
	for {
		select {
		case err := <-errChan:
			if err != nil {
				logrus.Errorf("error from request :: %v", err)
			}
		case sdp := <-sdpChan:
			handler.SDPRequest(sdp, errChan)
		}
	}
}
