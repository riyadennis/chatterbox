package handler

import (
	"fmt"
	"github.com/pion/webrtc/v2"
	"github.com/riyadennis/chatterbox/broadcast/signal"
	"github.com/riyadennis/chatterbox/internal"
	"net/http"
)

var (
	Port = "8087"
)

func SDPRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", fmt.Sprintf("localhost:%s", Port))
	offer, err := internal.OfferFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("cannot create offer"))
		return
	}
	peerCon, err := internal.PeerConnection(offer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("cannot create peer connection"))
		return
	}
	// Allow us to receive 1 video track
	if _, err = peerCon.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("cannot receive video track"))
		return
	}
	var localTrackChan chan *webrtc.Track
	var errChan chan error
	// Set a handler for when a new remote track starts, this just distributes all our packets
	// to connected peers
	peerCon.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
		localTrackChan, errChan = internal.TrackHandler(peerCon, remoteTrack)
		if err := <-errChan; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error creating local track"))
			return
		}
	})
	// Set the remote SessionDescription
	err = peerCon.SetRemoteDescription(offer)
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

	// Get the LocalDescription and take it to base64 so we can paste in browser
	fmt.Println(signal.Encode(answer))
	localTrack := <-localTrackChan
	for {
		fmt.Println("")
		fmt.Println("Curl an base64 SDP to start sendonly peer connection")

		recvOnlyOffer, err := internal.OfferFromRequest(r)
		if err != nil {
			panic(err)
		}
		// Create a new PeerConnection
		peerCon, err := internal.PeerConnection(offer)
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

		// Get the LocalDescription and take it to base64 so we can paste in browser
		fmt.Println(w, "%s", signal.Encode(answer))
	}
}
