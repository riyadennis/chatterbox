package internal

import (
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v2"
	"io"
	"time"
)

const (
	RtcpPLIInterval = time.Second * 3
)

func OfferFromRequest(body string) (webrtc.SessionDescription, error) {
	offer := webrtc.SessionDescription{}
	err := Decode(body, &offer)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}
	return offer, nil
}

func PeerConnection(offer webrtc.SessionDescription) (*webrtc.PeerConnection, error) {
	// Since we are answering use PayloadTypes declared by offerer
	mediaEngine := webrtc.MediaEngine{}
	err := mediaEngine.PopulateFromSDP(offer)
	if err != nil {
		return nil, err
	}
	// Create the API object with the MediaEngine
	api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))
	peerConnectionConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
	// Create a new RTCPeerConnection
	return api.NewPeerConnection(peerConnectionConfig)
}

func TrackHandler(peerCon *webrtc.PeerConnection, remoteTrack *webrtc.Track) (chan *webrtc.Track, chan error) {
	errChan := make(chan error)
	localTrackChan := make(chan *webrtc.Track)
	// Send a PLI on an interval so that the publisher is pushing a keyframe every RtcpPLIInterval
	// This can be less wasteful by processing incoming RTCP events, then we would emit a NACK/PLI when a viewer requests it
	go func() {
		ticker := time.NewTicker(RtcpPLIInterval)
		for range ticker.C {
			if rtcpSendErr := peerCon.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: remoteTrack.SSRC()}}); rtcpSendErr != nil {
				errChan <- rtcpSendErr
			}
		}
	}()
	// Create a local track, all our SFU clients will be fed via this track
	localTrack, newTrackErr := peerCon.NewTrack(remoteTrack.PayloadType(),
		remoteTrack.SSRC(), "video", "pion")
	if newTrackErr != nil {
		errChan <- newTrackErr
	}
	localTrackChan <- localTrack
	rtpBuf := make([]byte, 1400)
	for {
		i, readErr := remoteTrack.Read(rtpBuf)
		if readErr != nil {
			panic(readErr)
		}

		// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
		if _, err := localTrack.Write(rtpBuf[:i]); err != nil && err != io.ErrClosedPipe {
			errChan <- err
		}
	}
	return localTrackChan, errChan
}
