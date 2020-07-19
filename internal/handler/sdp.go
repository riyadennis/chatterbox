package handler

import (
	"fmt"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v2"
	"github.com/riyadennis/chatterbox/internal"
	"io"
	"time"
)

var (
	Port = "8080"
)

func SDPRequest(requestOffer string) error {
	offer, err := internal.OfferFromRequest(requestOffer)
	if err != nil {
		return err
	}
	peerCon, err := internal.PeerConnection(offer)
	if err != nil {
		return err
	}
	// Allow us to receive 1 video track
	if _, err = peerCon.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		return err
	}

	localTrackChan := make(chan *webrtc.Track)
	// Set a handler for when a new remote track starts, this just distributes all our packets
	// to connected peers
	peerCon.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
		// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
		// This can be less wasteful by processing incoming RTCP events, then we would emit a NACK/PLI when a viewer requests it
		go func() {
			ticker := time.NewTicker(internal.RtcpPLIInterval)
			for range ticker.C {
				if rtcpSendErr := peerCon.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: remoteTrack.SSRC()}}); rtcpSendErr != nil {
					fmt.Println(rtcpSendErr)
				}
			}
		}()

		// Create a local track, all our SFU clients will be fed via this track
		localTrack, newTrackErr := peerCon.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), "video", "pion")
		if newTrackErr != nil {
			panic(newTrackErr)
		}
		localTrackChan <- localTrack

		rtpBuf := make([]byte, 1400)
		for {
			i, readErr := remoteTrack.Read(rtpBuf)
			if readErr != nil {
				panic(readErr)
			}

			// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
			if _, err = localTrack.Write(rtpBuf[:i]); err != nil && err != io.ErrClosedPipe {
				panic(err)
			}
		}
	})

	// Set the remote SessionDescription
	err = peerCon.SetRemoteDescription(offer)
	if err != nil {
		return err
	}

	// Create answer
	answer, err := peerCon.CreateAnswer(nil)
	if err != nil {
		return err
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerCon.SetLocalDescription(answer)
	if err != nil {
		return err
	}

	// Get the LocalDescription and take it to base64 so we can paste in browser
	ans, err := internal.Encode(answer)
	if err != nil {
		return err
	}
	fmt.Printf("local desc outside for %s", ans)
	localTrack := <-localTrackChan
	for {
		fmt.Println("")
		fmt.Println("Curl an base64 SDP to start sendonly peer connection")

		recvOnlyOffer, err := internal.OfferFromRequest(requestOffer)
		if err != nil {
			panic(err)
		}
		// Create a new PeerConnection
		peerCon, err := internal.PeerConnection(recvOnlyOffer)
		if err != nil {
			panic(err)
		}

		_, err = peerCon.AddTrack(localTrack)
		if err != nil {
			panic(err)
		}

		// Set the remote SessionDescription
		err = peerCon.SetRemoteDescription(recvOnlyOffer)
		if err != nil {
			panic(err)
		}

		// Create answer
		answer, err := peerCon.CreateAnswer(nil)
		if err != nil {
			panic(err)
		}

		// Sets the LocalDescription, and starts our UDP listeners
		err = peerCon.SetLocalDescription(answer)
		if err != nil {
			panic(err)
		}
		ans, err := internal.Encode(answer)
		if err != nil {
			panic(err)
		}

		// Get the LocalDescription and take it to base64 so we can paste in browser
		fmt.Printf("local desc inside for %s", ans)
	}
}
