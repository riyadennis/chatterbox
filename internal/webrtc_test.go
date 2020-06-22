package internal

import (
	"github.com/pion/webrtc/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPeerConnection(t *testing.T) {
	tests := []struct {
		name          string
		sessionDesc   webrtc.SessionDescription
		expectedError string
	}{
		{
			name: "invalid offer type",
			sessionDesc: webrtc.SessionDescription{
				Type: 0,
				SDP:  "",
			},
			expectedError: "sdp: invalid syntax ``",
		},
		{
			name: "invalid sdp in offer",
			sessionDesc: func() webrtc.SessionDescription {
				return webrtc.SessionDescription{
					Type: 0,
					SDP:  "INVALID",
				}
			}(),
			expectedError: "sdp: invalid syntax `INVALID`",
		},
		{
			name: "invalid sdp in offer",
			sessionDesc: func() webrtc.SessionDescription {
				return webrtc.SessionDescription{
					Type: 0,
					SDP:  "v=0\r\no=- 1277610785989486581 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE 0\r\na=msid-semantic: WMS iNtS1DDBbo3KPy0Kt7noxNSGEwlOMavtuuZQ\r\nm=video 58780 UDP/TLS/RTP/SAVPF 96 97 98 99 100 101 102 122 127 121 125 107 108 109 124 120 123 119 114 115 116\r\nc=IN IP4 176.24.90.224\r\na=rtcp:9 IN IP4 0.0.0.0\r\na=candidate:2293670849 1 udp 2122262783 2a02:c7f:ac20:7800:a864:3aba:4c4a:92f4 51842 typ host generation 0 network-id 2 network-cost 10\r\na=candidate:1841357947 1 udp 2122194687 192.168.0.7 58780 typ host generation 0 network-id 1 network-cost 10\r\na=candidate:3298782126 1 udp 1685987071 176.24.90.224 58780 typ srflx raddr 192.168.0.7 rport 58780 generation 0 network-id 1 network-cost 10\r\na=candidate:3325386545 1 tcp 1518283007 2a02:c7f:ac20:7800:a864:3aba:4c4a:92f4 9 typ host tcptype active generation 0 network-id 2 network-cost 10\r\na=candidate:591599755 1 tcp 1518214911 192.168.0.7 9 typ host tcptype active generation 0 network-id 1 network-cost 10\r\na=ice-ufrag:F+xh\r\na=ice-pwd:ua/k9wxGVS+vtgIluYZBZEaW\r\na=ice-options:trickle\r\na=fingerprint:sha-256 6D:C9:EB:7D:D2:A9:C9:9D:46:19:0B:6B:1C:93:09:93:C5:B9:B3:F9:78:5C:B1:D0:24:A1:A4:BB:4D:D7:7F:31\r\na=setup:actpass\r\na=mid:0\r\na=extmap:14 urn:ietf:params:rtp-hdrext:toffset\r\na=extmap:2 http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time\r\na=extmap:13 urn:3gpp:video-orientation\r\na=extmap:3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01\r\na=extmap:12 http://www.webrtc.org/experiments/rtp-hdrext/playout-delay\r\na=extmap:11 http://www.webrtc.org/experiments/rtp-hdrext/video-content-type\r\na=extmap:7 http://www.webrtc.org/experiments/rtp-hdrext/video-timing\r\na=extmap:8 http://tools.ietf.org/html/draft-ietf-avtext-framemarking-07\r\na=extmap:9 http://www.webrtc.org/experiments/rtp-hdrext/color-space\r\na=extmap:4 urn:ietf:params:rtp-hdrext:sdes:mid\r\na=extmap:5 urn:ietf:params:rtp-hdrext:sdes:rtp-stream-id\r\na=extmap:6 urn:ietf:params:rtp-hdrext:sdes:repaired-rtp-stream-id\r\na=sendrecv\r\na=msid:iNtS1DDBbo3KPy0Kt7noxNSGEwlOMavtuuZQ ecc4f8b4-d595-4cd6-baa8-073f78084255\r\na=rtcp-mux\r\na=rtcp-rsize\r\na=rtpmap:96 VP8/90000\r\na=rtcp-fb:96 goog-remb\r\na=rtcp-fb:96 transport-cc\r\na=rtcp-fb:96 ccm fir\r\na=rtcp-fb:96 nack\r\na=rtcp-fb:96 nack pli\r\na=rtpmap:97 rtx/90000\r\na=fmtp:97 apt=96\r\na=rtpmap:98 VP9/90000\r\na=rtcp-fb:98 goog-remb\r\na=rtcp-fb:98 transport-cc\r\na=rtcp-fb:98 ccm fir\r\na=rtcp-fb:98 nack\r\na=rtcp-fb:98 nack pli\r\na=fmtp:98 profile-id=0\r\na=rtpmap:99 rtx/90000\r\na=fmtp:99 apt=98\r\na=rtpmap:100 VP9/90000\r\na=rtcp-fb:100 goog-remb\r\na=rtcp-fb:100 transport-cc\r\na=rtcp-fb:100 ccm fir\r\na=rtcp-fb:100 nack\r\na=rtcp-fb:100 nack pli\r\na=fmtp:100 profile-id=2\r\na=rtpmap:101 rtx/90000\r\na=fmtp:101 apt=100\r\na=rtpmap:102 H264/90000\r\na=rtcp-fb:102 goog-remb\r\na=rtcp-fb:102 transport-cc\r\na=rtcp-fb:102 ccm fir\r\na=rtcp-fb:102 nack\r\na=rtcp-fb:102 nack pli\r\na=fmtp:102 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42001f\r\na=rtpmap:122 rtx/90000\r\na=fmtp:122 apt=102\r\na=rtpmap:127 H264/90000\r\na=rtcp-fb:127 goog-remb\r\na=rtcp-fb:127 transport-cc\r\na=rtcp-fb:127 ccm fir\r\na=rtcp-fb:127 nack\r\na=rtcp-fb:127 nack pli\r\na=fmtp:127 level-asymmetry-allowed=1;packetization-mode=0;profile-level-id=42001f\r\na=rtpmap:121 rtx/90000\r\na=fmtp:121 apt=127\r\na=rtpmap:125 H264/90000\r\na=rtcp-fb:125 goog-remb\r\na=rtcp-fb:125 transport-cc\r\na=rtcp-fb:125 ccm fir\r\na=rtcp-fb:125 nack\r\na=rtcp-fb:125 nack pli\r\na=fmtp:125 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42e01f\r\na=rtpmap:107 rtx/90000\r\na=fmtp:107 apt=125\r\na=rtpmap:108 H264/90000\r\na=rtcp-fb:108 goog-remb\r\na=rtcp-fb:108 transport-cc\r\na=rtcp-fb:108 ccm fir\r\na=rtcp-fb:108 nack\r\na=rtcp-fb:108 nack pli\r\na=fmtp:108 level-asymmetry-allowed=1;packetization-mode=0;profile-level-id=42e01f\r\na=rtpmap:109 rtx/90000\r\na=fmtp:109 apt=108\r\na=rtpmap:124 H264/90000\r\na=rtcp-fb:124 goog-remb\r\na=rtcp-fb:124 transport-cc\r\na=rtcp-fb:124 ccm fir\r\na=rtcp-fb:124 nack\r\na=rtcp-fb:124 nack pli\r\na=fmtp:124 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=4d0032\r\na=rtpmap:120 rtx/90000\r\na=fmtp:120 apt=124\r\na=rtpmap:123 H264/90000\r\na=rtcp-fb:123 goog-remb\r\na=rtcp-fb:123 transport-cc\r\na=rtcp-fb:123 ccm fir\r\na=rtcp-fb:123 nack\r\na=rtcp-fb:123 nack pli\r\na=fmtp:123 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=640032\r\na=rtpmap:119 rtx/90000\r\na=fmtp:119 apt=123\r\na=rtpmap:114 red/90000\r\na=rtpmap:115 rtx/90000\r\na=fmtp:115 apt=114\r\na=rtpmap:116 ulpfec/90000\r\na=ssrc-group:FID 1873160381 3829752818\r\na=ssrc:1873160381 cname:ZgITmDEfY7pyEwXM\r\na=ssrc:1873160381 msid:iNtS1DDBbo3KPy0Kt7noxNSGEwlOMavtuuZQ ecc4f8b4-d595-4cd6-baa8-073f78084255\r\na=ssrc:1873160381 mslabel:iNtS1DDBbo3KPy0Kt7noxNSGEwlOMavtuuZQ\r\na=ssrc:1873160381 label:ecc4f8b4-d595-4cd6-baa8-073f78084255\r\na=ssrc:3829752818 cname:ZgITmDEfY7pyEwXM\r\na=ssrc:3829752818 msid:iNtS1DDBbo3KPy0Kt7noxNSGEwlOMavtuuZQ ecc4f8b4-d595-4cd6-baa8-073f78084255\r\na=ssrc:3829752818 mslabel:iNtS1DDBbo3KPy0Kt7noxNSGEwlOMavtuuZQ\r\na=ssrc:3829752818 label:ecc4f8b4-d595-4cd6-baa8-073f78084255\r\n",
				}
			}(),
			expectedError: "sdp: invalid syntax `INVALID`",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			conn, err := PeerConnection(tc.sessionDesc)
			if err != nil {
				assert.Equal(t, tc.expectedError, err.Error())
				return
			}
			assert.NotNil(t, conn, "empty connection")
		})
	}
}
